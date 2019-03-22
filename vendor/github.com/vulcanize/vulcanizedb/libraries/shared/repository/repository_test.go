// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package repository_test

import (
	"fmt"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	shared "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	r2 "github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Repository", func() {
	Describe("MarkHeaderChecked", func() {
		var (
			checkedHeadersColumn string
			db                   *postgres.DB
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers
				ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		It("marks passed column as checked for passed header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())

			err := shared.MarkHeaderChecked(headerID, db, checkedHeadersColumn)

			Expect(err).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT `+checkedHeadersColumn+` FROM public.checked_headers LIMIT 1`)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})
	})

	Describe("MarkHeaderCheckedInTransaction", func() {
		var (
			checkedHeadersColumn string
			db                   *postgres.DB
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers
				ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		It("marks passed column as checked for passed header within a passed transaction", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			err := shared.MarkHeaderCheckedInTransaction(headerID, tx, checkedHeadersColumn)

			Expect(err).NotTo(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT `+checkedHeadersColumn+` FROM public.checked_headers LIMIT 1`)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})
	})

	Describe("MissingHeaders", func() {
		var (
			db                       *postgres.DB
			headerRepository         datastore.HeaderRepository
			startingBlockNumber      int64
			endingBlockNumber        int64
			eventSpecificBlockNumber int64
			blockNumbers             []int64
			headerIDs                []int64
			notCheckedSQL            string
			err                      error
			hr                       r2.HeaderRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			headerRepository = repositories.NewHeaderRepository(db)
			hr = r2.NewHeaderRepository(db)
			hr.AddCheckColumns(getExpectedColumnNames())

			columnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			notCheckedSQL = shared.CreateNotCheckedSQL(columnNames, constants.HeaderMissing)

			startingBlockNumber = rand.Int63()
			eventSpecificBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2
			outOfRangeBlockNumber := endingBlockNumber + 1

			blockNumbers = []int64{startingBlockNumber, eventSpecificBlockNumber, endingBlockNumber, outOfRangeBlockNumber}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			test_config.CleanCheckedHeadersTable(db, getExpectedColumnNames())
		})

		It("only treats headers as checked if the event specific logs have been checked", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n + 10))
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(err).NotTo(HaveOccurred())
			nodeOneMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(3))
			Expect(nodeOneMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))

			nodeTwoMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber+10, dbTwo, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
			Expect(nodeTwoMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
			Expect(nodeTwoMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
		})
	})

	Describe("RecheckHeaders", func() {
		var (
			checkedHeadersColumn                                      string
			db                                                        *postgres.DB
			headerOneID, headerTwoID, headerThreeID, headerFourID     int64
			headerOneErr, headerTwoErr, headerThreeErr, headerFourErr error
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)

			// create header checked column
			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())

			// create headers
			headerRepository := repositories.NewHeaderRepository(db)
			headerOneID, headerOneErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(1))
			Expect(headerOneErr).NotTo(HaveOccurred())
			headerTwoID, headerTwoErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(2))
			Expect(headerTwoErr).NotTo(HaveOccurred())
			headerThreeID, headerThreeErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(3))
			Expect(headerThreeErr).NotTo(HaveOccurred())
			headerFourID, headerFourErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(4))
			Expect(headerFourErr).NotTo(HaveOccurred())

			// mark every header checked at least once, with one fully rechecked (headerThree)
			maxCheckCount, intConversionErr := strconv.Atoi(constants.RecheckHeaderCap)
			Expect(intConversionErr).NotTo(HaveOccurred())
			_, markHeaderOneCheckedErr := db.Exec(
				`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`) VALUES ($1, $2)`,
				headerOneID, maxCheckCount)
			Expect(markHeaderOneCheckedErr).NotTo(HaveOccurred())
			_, markHeaderTwoCheckedErr := db.Exec(
				`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`) VALUES ($1, $2)`,
				headerTwoID, maxCheckCount)
			Expect(markHeaderTwoCheckedErr).NotTo(HaveOccurred())
			_, markHeaderThreeCheckedErr := db.Exec(
				`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`) VALUES ($1, $2)`,
				headerThreeID, maxCheckCount+1)
			Expect(markHeaderThreeCheckedErr).NotTo(HaveOccurred())
			_, markHeaderFourCheckedErr := db.Exec(
				`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`) VALUES ($1, $2)`,
				headerFourID, maxCheckCount)
			Expect(markHeaderFourCheckedErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		Describe("when no ending block number (ending block number == -1)", func() {
			It("returns all headers since starting block where checked count is less than cap", func() {
				headers, err := shared.RecheckHeaders(1, -1, db, checkedHeadersColumn)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(headers)).To(Equal(3))
				Expect(headers[0].Id).To(Or(Equal(headerOneID), Equal(headerTwoID), Equal(headerFourID)))
				Expect(headers[1].Id).To(Or(Equal(headerOneID), Equal(headerTwoID), Equal(headerFourID)))
				Expect(headers[2].Id).To(Or(Equal(headerOneID), Equal(headerTwoID), Equal(headerFourID)))
			})
		})

		Describe("when ending block number specified", func() {
			It("returns headers between starting and ending block where checked count is less than cap", func() {
				headers, err := shared.RecheckHeaders(1, 3, db, checkedHeadersColumn)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(headers)).To(Equal(2))
				Expect(headers[0].Id).To(Or(Equal(headerOneID), Equal(headerTwoID)))
				Expect(headers[1].Id).To(Or(Equal(headerOneID), Equal(headerTwoID)))
			})
		})
	})

	Describe("GetCheckedColumnNames", func() {
		It("gets the column names from checked_headers", func() {
			db := test_config.NewTestDB(test_config.NewTestNode())
			hr := r2.NewHeaderRepository(db)
			hr.AddCheckColumns(getExpectedColumnNames())
			test_config.CleanTestDB(db)
			expectedColumnNames := getExpectedColumnNames()
			actualColumnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualColumnNames).To(Equal(expectedColumnNames))
			test_config.CleanCheckedHeadersTable(db, getExpectedColumnNames())
		})
	})

	Describe("CreateNotCheckedSQL", func() {
		It("generates a correct SQL string for one column", func() {
			columns := []string{"columnA"}
			expected := "NOT (columnA!=0)"
			actual := shared.CreateNotCheckedSQL(columns, constants.HeaderMissing)
			Expect(actual).To(Equal(expected))
		})

		It("generates a correct SQL string for several columns", func() {
			columns := []string{"columnA", "columnB"}
			expected := "NOT (columnA!=0 AND columnB!=0)"
			actual := shared.CreateNotCheckedSQL(columns, constants.HeaderMissing)
			Expect(actual).To(Equal(expected))
		})

		It("defaults to FALSE when there are no columns", func() {
			expected := "FALSE"
			actual := shared.CreateNotCheckedSQL([]string{}, constants.HeaderMissing)
			Expect(actual).To(Equal(expected))
		})

		It("generates a correct SQL string for rechecking headers", func() {
			columns := []string{"columnA", "columnB"}
			expected := fmt.Sprintf("NOT (columnA>=%s AND columnB>=%s)", constants.RecheckHeaderCap, constants.RecheckHeaderCap)
			actual := shared.CreateNotCheckedSQL(columns, constants.HeaderRecheck)
			Expect(actual).To(Equal(expected))
		})
	})
})

func getExpectedColumnNames() []string {
	return []string{
		"column_1_checked",
		"column_2_checked",
		"column_3_checked",
		"column_4_checked",
	}
}
