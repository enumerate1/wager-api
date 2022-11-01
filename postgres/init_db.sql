-- wager table
CREATE SEQUENCE IF NOT EXISTS public.wager_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE IF NOT EXISTS public.wager (
    wager_id integer NOT NULL DEFAULT nextval('wager_id_seq'),
    total_wager_value REAL NOT NULL,
    odds INTEGER,
    selling_percentage INTEGER,
    selling_price REAL,
    current_selling_price REAL,
    percentage_sold REAL,
    amount_sold REAL,
    place_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT wager_pk PRIMARY KEY (wager_id)

);
ALTER SEQUENCE IF EXISTS wager_id_seq OWNED BY wager.wager_id;

-- purchase table
CREATE SEQUENCE IF NOT EXISTS public.purchase_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE IF NOT EXISTS public.purchase (
    purchase_id integer NOT NULL DEFAULT nextval('purchase_id_seq'),
    wager_id integer,
    buying_price REAL,
    bought_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT purchase_pk PRIMARY KEY (purchase_id),
    CONSTRAINT purchase_wager_fk FOREIGN KEY (wager_id) REFERENCES public.wager(wager_id)
);
ALTER SEQUENCE IF EXISTS purchase_id_seq OWNED BY purchase.purchase_id;
