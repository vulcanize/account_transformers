--
-- PostgreSQL database dump
--

-- Dumped from database version 10.6 (Ubuntu 10.6-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.7 (Ubuntu 10.7-1.pgdg16.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: accounts; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA accounts;


--
-- Name: postgraphile_watch; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA postgraphile_watch;


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: notify_watchers_ddl(); Type: FUNCTION; Schema: postgraphile_watch; Owner: -
--

CREATE FUNCTION postgraphile_watch.notify_watchers_ddl() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(
    'postgraphile_watch',
    json_build_object(
      'type',
      'ddl',
      'payload',
      (select json_agg(json_build_object('schema', schema_name, 'command', command_tag)) from pg_event_trigger_ddl_commands() as x)
    )::text
  );
end;
$$;


--
-- Name: notify_watchers_drop(); Type: FUNCTION; Schema: postgraphile_watch; Owner: -
--

CREATE FUNCTION postgraphile_watch.notify_watchers_drop() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(
    'postgraphile_watch',
    json_build_object(
      'type',
      'drop',
      'payload',
      (select json_agg(distinct x.schema_name) from pg_event_trigger_dropped_objects() as x)
    )::text
  );
end;
$$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: address_coin_balances; Type: TABLE; Schema: accounts; Owner: -
--

CREATE TABLE accounts.address_coin_balances (
    id integer NOT NULL,
    header_id integer NOT NULL,
    address_hash character varying(42) NOT NULL,
    block_number bigint NOT NULL,
    value numeric(100,0),
    value_fetched_at timestamp without time zone,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone
);


--
-- Name: address_coin_balances_id_seq; Type: SEQUENCE; Schema: accounts; Owner: -
--

CREATE SEQUENCE accounts.address_coin_balances_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: address_coin_balances_id_seq; Type: SEQUENCE OWNED BY; Schema: accounts; Owner: -
--

ALTER SEQUENCE accounts.address_coin_balances_id_seq OWNED BY accounts.address_coin_balances.id;


--
-- Name: addresses; Type: TABLE; Schema: accounts; Owner: -
--

CREATE TABLE accounts.addresses (
    address character varying(42) NOT NULL
);


--
-- Name: contract_addresses; Type: TABLE; Schema: accounts; Owner: -
--

CREATE TABLE accounts.contract_addresses (
    contract character varying(42) NOT NULL
);


--
-- Name: token_value_transfers; Type: TABLE; Schema: accounts; Owner: -
--

CREATE TABLE accounts.token_value_transfers (
    id integer NOT NULL,
    header_id integer NOT NULL,
    block_number bigint NOT NULL,
    name character varying NOT NULL,
    dst character varying(42),
    src character varying(42),
    amount numeric,
    contract character varying(42) NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb,
    CONSTRAINT token_value_transfers_name_check CHECK (((name)::text <> ''::text))
);


--
-- Name: headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66),
    block_number bigint,
    raw jsonb,
    block_timestamp numeric,
    eth_node_id integer,
    eth_node_fingerprint character varying(128)
);


--
-- Name: address_token_balances; Type: VIEW; Schema: accounts; Owner: -
--

CREATE VIEW accounts.address_token_balances AS
 SELECT addresses.address AS address_hash,
    contract_addresses.contract AS token_contract_address_hash,
    headers.block_number,
    (( SELECT COALESCE(sum(token_value_transfers_1.amount), (0)::numeric) AS "coalesce"
           FROM accounts.token_value_transfers token_value_transfers_1
          WHERE ((token_value_transfers_1.block_number <= headers.block_number) AND ((token_value_transfers_1.dst)::text = (addresses.address)::text) AND ((token_value_transfers_1.contract)::text = (contract_addresses.contract)::text))) - ( SELECT COALESCE(sum(token_value_transfers_1.amount), (0)::numeric) AS "coalesce"
           FROM accounts.token_value_transfers token_value_transfers_1
          WHERE ((token_value_transfers_1.block_number <= headers.block_number) AND ((token_value_transfers_1.src)::text = (addresses.address)::text) AND ((token_value_transfers_1.contract)::text = (contract_addresses.contract)::text)))) AS value
   FROM accounts.token_value_transfers,
    accounts.addresses,
    public.headers,
    accounts.contract_addresses
  GROUP BY addresses.address, contract_addresses.contract, headers.block_number;


--
-- Name: token_value_transfers_id_seq; Type: SEQUENCE; Schema: accounts; Owner: -
--

CREATE SEQUENCE accounts.token_value_transfers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: token_value_transfers_id_seq; Type: SEQUENCE OWNED BY; Schema: accounts; Owner: -
--

ALTER SEQUENCE accounts.token_value_transfers_id_seq OWNED BY accounts.token_value_transfers.id;


--
-- Name: logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logs (
    id integer NOT NULL,
    block_number bigint,
    address character varying(66),
    tx_hash character varying(66),
    index bigint,
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    data text,
    receipt_id integer
);


--
-- Name: block_stats; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.block_stats AS
 SELECT max(logs.block_number) AS max_block,
    min(logs.block_number) AS min_block
   FROM public.logs;


--
-- Name: blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blocks (
    id integer NOT NULL,
    difficulty bigint,
    extra_data character varying,
    gaslimit bigint,
    gasused bigint,
    hash character varying(66),
    miner character varying(42),
    nonce character varying(20),
    number bigint,
    parenthash character varying(66),
    reward numeric,
    uncles_reward numeric,
    size character varying,
    "time" bigint,
    is_final boolean,
    uncle_hash character varying(66),
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128) NOT NULL
);


--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


--
-- Name: checked_headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.checked_headers (
    id integer NOT NULL,
    header_id integer NOT NULL,
    token_value_transfers integer DEFAULT 0 NOT NULL
);


--
-- Name: checked_headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.checked_headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: checked_headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.checked_headers_id_seq OWNED BY public.checked_headers.id;


--
-- Name: eth_nodes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.eth_nodes (
    id integer NOT NULL,
    client_name character varying,
    genesis_block character varying(66),
    network_id numeric,
    eth_node_id character varying(128)
);


--
-- Name: full_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_receipts (
    id integer NOT NULL,
    contract_address character varying(42),
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    block_id integer NOT NULL
);


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_receipts_id_seq OWNED BY public.full_sync_receipts.id;


--
-- Name: full_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_transactions (
    id integer NOT NULL,
    block_id integer NOT NULL,
    gaslimit numeric,
    gasprice numeric,
    hash character varying(66),
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(66),
    tx_index integer,
    tx_to character varying(66),
    value numeric
);


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_transactions_id_seq OWNED BY public.full_sync_transactions.id;


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.headers_id_seq OWNED BY public.headers.id;


--
-- Name: light_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.light_sync_receipts (
    id integer NOT NULL,
    transaction_id integer NOT NULL,
    header_id integer NOT NULL,
    contract_address character varying(42),
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66)
);


--
-- Name: light_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.light_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: light_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.light_sync_receipts_id_seq OWNED BY public.light_sync_receipts.id;


--
-- Name: light_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.light_sync_transactions (
    id integer NOT NULL,
    header_id integer NOT NULL,
    hash character varying(66),
    gaslimit numeric,
    gasprice numeric,
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(44),
    tx_index integer,
    tx_to character varying(44),
    value numeric
);


--
-- Name: light_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.light_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: light_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.light_sync_transactions_id_seq OWNED BY public.light_sync_transactions.id;


--
-- Name: log_filters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.log_filters (
    id integer NOT NULL,
    name character varying NOT NULL,
    from_block bigint,
    to_block bigint,
    address character varying(66),
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    CONSTRAINT log_filters_from_block_check CHECK ((from_block >= 0)),
    CONSTRAINT log_filters_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT log_filters_to_block_check CHECK ((to_block >= 0))
);


--
-- Name: log_filters_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.log_filters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.log_filters_id_seq OWNED BY public.log_filters.id;


--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


--
-- Name: nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: queued_storage; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.queued_storage (
    id integer NOT NULL,
    block_height bigint,
    block_hash bytea,
    contract bytea,
    storage_key bytea,
    storage_value bytea
);


--
-- Name: queued_storage_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.queued_storage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: queued_storage_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.queued_storage_id_seq OWNED BY public.queued_storage.id;


--
-- Name: uncles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uncles (
    id integer NOT NULL,
    hash character varying(66) NOT NULL,
    block_id integer NOT NULL,
    block_hash character varying(66) NOT NULL,
    reward numeric NOT NULL,
    miner character varying(42) NOT NULL,
    raw jsonb,
    block_timestamp numeric,
    eth_node_id integer,
    eth_node_fingerprint character varying(128)
);


--
-- Name: uncles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.uncles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uncles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.uncles_id_seq OWNED BY public.uncles.id;


--
-- Name: watched_contracts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.watched_contracts (
    contract_id integer NOT NULL,
    contract_abi json,
    contract_hash character varying(66)
);


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.watched_contracts_contract_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.watched_contracts_contract_id_seq OWNED BY public.watched_contracts.contract_id;


--
-- Name: watched_event_logs; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.watched_event_logs AS
 SELECT log_filters.name,
    logs.id,
    logs.block_number,
    logs.address,
    logs.tx_hash,
    logs.index,
    logs.topic0,
    logs.topic1,
    logs.topic2,
    logs.topic3,
    logs.data,
    logs.receipt_id
   FROM ((public.log_filters
     CROSS JOIN public.block_stats)
     JOIN public.logs ON ((((logs.address)::text = (log_filters.address)::text) AND (logs.block_number >= COALESCE(log_filters.from_block, block_stats.min_block)) AND (logs.block_number <= COALESCE(log_filters.to_block, block_stats.max_block)))))
  WHERE ((((log_filters.topic0)::text = (logs.topic0)::text) OR (log_filters.topic0 IS NULL)) AND (((log_filters.topic1)::text = (logs.topic1)::text) OR (log_filters.topic1 IS NULL)) AND (((log_filters.topic2)::text = (logs.topic2)::text) OR (log_filters.topic2 IS NULL)) AND (((log_filters.topic3)::text = (logs.topic3)::text) OR (log_filters.topic3 IS NULL)));


--
-- Name: address_coin_balances id; Type: DEFAULT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.address_coin_balances ALTER COLUMN id SET DEFAULT nextval('accounts.address_coin_balances_id_seq'::regclass);


--
-- Name: token_value_transfers id; Type: DEFAULT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.token_value_transfers ALTER COLUMN id SET DEFAULT nextval('accounts.token_value_transfers_id_seq'::regclass);


--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);


--
-- Name: full_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.full_sync_receipts_id_seq'::regclass);


--
-- Name: full_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.full_sync_transactions_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: light_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.light_sync_receipts_id_seq'::regclass);


--
-- Name: light_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.light_sync_transactions_id_seq'::regclass);


--
-- Name: log_filters id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters ALTER COLUMN id SET DEFAULT nextval('public.log_filters_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: queued_storage id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage ALTER COLUMN id SET DEFAULT nextval('public.queued_storage_id_seq'::regclass);


--
-- Name: uncles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles ALTER COLUMN id SET DEFAULT nextval('public.uncles_id_seq'::regclass);


--
-- Name: watched_contracts contract_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts ALTER COLUMN contract_id SET DEFAULT nextval('public.watched_contracts_contract_id_seq'::regclass);


--
-- Name: address_coin_balances address_coin_balances_address_hash_block_number_key; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.address_coin_balances
    ADD CONSTRAINT address_coin_balances_address_hash_block_number_key UNIQUE (address_hash, block_number);


--
-- Name: address_coin_balances address_coin_balances_pkey; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.address_coin_balances
    ADD CONSTRAINT address_coin_balances_pkey PRIMARY KEY (id);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (address);


--
-- Name: contract_addresses contract_addresses_pkey; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.contract_addresses
    ADD CONSTRAINT contract_addresses_pkey PRIMARY KEY (contract);


--
-- Name: token_value_transfers token_value_transfers_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.token_value_transfers
    ADD CONSTRAINT token_value_transfers_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: token_value_transfers token_value_transfers_pkey; Type: CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.token_value_transfers
    ADD CONSTRAINT token_value_transfers_pkey PRIMARY KEY (id);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


--
-- Name: checked_headers checked_headers_header_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_key UNIQUE (header_id);


--
-- Name: checked_headers checked_headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_pkey PRIMARY KEY (id);


--
-- Name: blocks eth_node_id_block_number_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);


--
-- Name: eth_nodes eth_node_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);


--
-- Name: full_sync_receipts full_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT full_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: full_sync_transactions full_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: light_sync_receipts light_sync_receipts_header_id_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_receipts
    ADD CONSTRAINT light_sync_receipts_header_id_transaction_id_key UNIQUE (header_id, transaction_id);


--
-- Name: light_sync_receipts light_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_receipts
    ADD CONSTRAINT light_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: light_sync_transactions light_sync_transactions_header_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_header_id_hash_key UNIQUE (header_id, hash);


--
-- Name: light_sync_transactions light_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: log_filters name_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters
    ADD CONSTRAINT name_uc UNIQUE (name);


--
-- Name: eth_nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


--
-- Name: queued_storage queued_storage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_pkey PRIMARY KEY (id);


--
-- Name: uncles uncles_block_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_hash_key UNIQUE (block_id, hash);


--
-- Name: uncles uncles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_pkey PRIMARY KEY (id);


--
-- Name: watched_contracts watched_contracts_contract_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_contract_hash_key UNIQUE (contract_hash);


--
-- Name: watched_contracts watched_contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_pkey PRIMARY KEY (contract_id);


--
-- Name: address_coin_balances_address_hash_block_number_index; Type: INDEX; Schema: accounts; Owner: -
--

CREATE UNIQUE INDEX address_coin_balances_address_hash_block_number_index ON accounts.address_coin_balances USING btree (address_hash, block_number);


--
-- Name: unfetched_balances; Type: INDEX; Schema: accounts; Owner: -
--

CREATE UNIQUE INDEX unfetched_balances ON accounts.address_coin_balances USING btree (address_hash, block_number) WHERE (value_fetched_at IS NULL);


--
-- Name: block_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX block_id_index ON public.full_sync_transactions USING btree (block_id);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: node_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX node_id_index ON public.blocks USING btree (eth_node_id);


--
-- Name: number_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX number_index ON public.blocks USING btree (number);


--
-- Name: tx_from_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_from_index ON public.full_sync_transactions USING btree (tx_from);


--
-- Name: tx_to_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_to_index ON public.full_sync_transactions USING btree (tx_to);


--
-- Name: address_coin_balances address_coin_balances_header_id_fkey; Type: FK CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.address_coin_balances
    ADD CONSTRAINT address_coin_balances_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: token_value_transfers token_value_transfers_header_id_fkey; Type: FK CONSTRAINT; Schema: accounts; Owner: -
--

ALTER TABLE ONLY accounts.token_value_transfers
    ADD CONSTRAINT token_value_transfers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: full_sync_receipts blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: headers eth_nodes_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: uncles eth_nodes_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: full_sync_transactions full_sync_transactions_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: light_sync_receipts light_sync_receipts_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_receipts
    ADD CONSTRAINT light_sync_receipts_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: light_sync_receipts light_sync_receipts_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_receipts
    ADD CONSTRAINT light_sync_receipts_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.light_sync_transactions(id) ON DELETE CASCADE;


--
-- Name: light_sync_transactions light_sync_transactions_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: blocks node_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT node_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT receipts_fk FOREIGN KEY (receipt_id) REFERENCES public.full_sync_receipts(id) ON DELETE CASCADE;


--
-- Name: uncles uncles_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: postgraphile_watch_ddl; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER postgraphile_watch_ddl ON ddl_command_end
         WHEN TAG IN ('ALTER AGGREGATE', 'ALTER DOMAIN', 'ALTER EXTENSION', 'ALTER FOREIGN TABLE', 'ALTER FUNCTION', 'ALTER POLICY', 'ALTER SCHEMA', 'ALTER TABLE', 'ALTER TYPE', 'ALTER VIEW', 'COMMENT', 'CREATE AGGREGATE', 'CREATE DOMAIN', 'CREATE EXTENSION', 'CREATE FOREIGN TABLE', 'CREATE FUNCTION', 'CREATE INDEX', 'CREATE POLICY', 'CREATE RULE', 'CREATE SCHEMA', 'CREATE TABLE', 'CREATE TABLE AS', 'CREATE VIEW', 'DROP AGGREGATE', 'DROP DOMAIN', 'DROP EXTENSION', 'DROP FOREIGN TABLE', 'DROP FUNCTION', 'DROP INDEX', 'DROP OWNED', 'DROP POLICY', 'DROP RULE', 'DROP SCHEMA', 'DROP TABLE', 'DROP TYPE', 'DROP VIEW', 'GRANT', 'REVOKE', 'SELECT INTO')
   EXECUTE PROCEDURE postgraphile_watch.notify_watchers_ddl();


--
-- Name: postgraphile_watch_drop; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER postgraphile_watch_drop ON sql_drop
   EXECUTE PROCEDURE postgraphile_watch.notify_watchers_drop();


--
-- PostgreSQL database dump complete
--

