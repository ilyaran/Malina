--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.1
-- Dumped by pg_dump version 9.6.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: cart_product; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE cart_product AS (
	product_id bigint,
	product_img text,
	product_title character varying(255),
	product_price numeric(20,2),
	product_quantity bigint,
	subtotal numeric(20,2)
);


ALTER TYPE cart_product OWNER TO postgres;

--
-- Name: cart_saved_product; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE cart_saved_product AS (
	product_id bigint,
	product_title character varying(255),
	product_price numeric(20,2)
);


ALTER TYPE cart_saved_product OWNER TO postgres;

--
-- Name: cart_add_product(character, bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_add_product(character, bigint) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  inProductId ALIAS FOR $2;
  productQuantity INTEGER;
BEGIN
  SELECT INTO productQuantity
    cart_quantity
  FROM cart
  WHERE cart_id = inCartId AND cart_product = inProductId;
  IF productQuantity IS NULL THEN
    INSERT INTO cart(cart_id, cart_product, cart_quantity, cart_created)
    VALUES (inCartId, inProductId , 1, NOW());
  ELSE
    UPDATE cart
    SET cart_quantity = cart_quantity + 1, cart_buy_now = true
    WHERE cart_id = inCartId AND cart_product = inProductId;
  END IF;
END;
$_$;


ALTER FUNCTION public.cart_add_product(character, bigint) OWNER TO postgres;

--
-- Name: cart_get_products(character); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_get_products(character) RETURNS SETOF cart_product
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  outCartProductRow cart_product;
BEGIN
  FOR outCartProductRow IN
  SELECT p.product_id, array_to_string(p.product_img,'|'), p.product_title,
    COALESCE(NULLIF(p.product_price1, 0), p.product_price) AS price,
    sc.cart_quantity,
    COALESCE(NULLIF(p.product_price1, 0),
             p.product_price) * sc.cart_quantity AS subtotal
  FROM cart sc
    INNER JOIN product p
      ON sc.cart_product = p.product_id
  WHERE sc.cart_id = inCartId AND sc.cart_buy_now
  LOOP
    RETURN NEXT outCartProductRow;
  END LOOP;
END;
$_$;


ALTER FUNCTION public.cart_get_products(character) OWNER TO postgres;

--
-- Name: cart_get_saved_products(character); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_get_saved_products(character) RETURNS SETOF cart_saved_product
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  outCartSavedProductRow cart_saved_product;
BEGIN
  FOR outCartSavedProductRow IN
  SELECT p.product_id, p.product_title,
    COALESCE(NULLIF(p.product_price1, 0), p.product_price) AS price
  FROM cart sc
    INNER JOIN product p
      ON sc.cart_product = p.product_id
  WHERE sc.cart_id = inCartId AND NOT cart_buy_now
  LOOP
  RETURN NEXT outCartSavedProductRow;
END LOOP;
END;
$_$;


ALTER FUNCTION public.cart_get_saved_products(character) OWNER TO postgres;

--
-- Name: cart_get_total_amount(character); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_get_total_amount(character) RETURNS numeric
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  outTotalAmount NUMERIC(20, 2);
BEGIN
  SELECT INTO outTotalAmount
    SUM(COALESCE(NULLIF(p.product_price1, 0), p.product_price)
        * sc.cart_quantity)
  FROM cart sc
    INNER JOIN product p
      ON sc.cart_product = p.product_id
  WHERE sc.cart_id = inCartId AND sc.cart_buy_now;
  RETURN outTotalAmount;
END;
$_$;


ALTER FUNCTION public.cart_get_total_amount(character) OWNER TO postgres;

--
-- Name: cart_move_product_to_cart(character, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_move_product_to_cart(character, integer) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  inProductId ALIAS FOR $2;
BEGIN
  UPDATE cart
  SET cart_buy_now = true, cart_created = NOW()
  WHERE cart_id = inCartId AND cart_product = inProductId;
END;
$_$;


ALTER FUNCTION public.cart_move_product_to_cart(character, integer) OWNER TO postgres;

--
-- Name: cart_remove_product(character, bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_remove_product(character, bigint) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE
inCartId ALIAS FOR $1;
inProductId ALIAS FOR $2;
BEGIN
DELETE FROM cart
WHERE cart_id = inCartId AND cart_product = inProductId;
END;
$_$;


ALTER FUNCTION public.cart_remove_product(character, bigint) OWNER TO postgres;

--
-- Name: cart_save_product_for_later(character, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_save_product_for_later(character, integer) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  inProductId ALIAS FOR $2;
BEGIN
  UPDATE cart
  SET cart_buy_now = false, cart_quantity = 1
  WHERE cart_id = inCartId AND cart_product = inProductId;
END;
$_$;


ALTER FUNCTION public.cart_save_product_for_later(character, integer) OWNER TO postgres;

--
-- Name: cart_update(character, bigint[], bigint[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION cart_update(character, bigint[], bigint[]) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCartId ALIAS FOR $1;
  inProductIds ALIAS FOR $2;
  inQuantities ALIAS FOR $3;
BEGIN
  FOR i IN array_lower(inQuantities, 1)..array_upper(inQuantities, 1)
  LOOP
    IF inQuantities[i] > 0 THEN
      UPDATE cart
      SET cart_quantity = inQuantities[i], cart_created = NOW()
      WHERE cart_id = inCartId AND cart_product = inProductIds[i];
    ELSE
      PERFORM cart_remove_product(inCartId, inProductIds[i]);
    END IF;
  END LOOP;
END;
$_$;


ALTER FUNCTION public.cart_update(character, bigint[], bigint[]) OWNER TO postgres;

--
-- Name: product_add(bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION product_add(bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]) RETURNS bigint
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inCategoryId      ALIAS FOR $1;
  inTitle           ALIAS FOR $2;
  inDescription     ALIAS FOR $3;
  inCode            ALIAS FOR $4;
  inPrice           ALIAS FOR $5;
  inPrice1          ALIAS FOR $6;
  inEnable          ALIAS FOR $7;
  inImgId           ALIAS FOR $8;
  lastInsertId BIGINT;
BEGIN

  INSERT INTO product (
    product_title,
    product_description,
    product_code,
    product_price,
    product_price1,
    product_enable,
    product_img,
    search_vector)
  VALUES
    (inTitle, inDescription,inCode,inPrice,inPrice1,inEnable,inImgId,
     (setweight(to_tsvector(inTitle), 'A') || to_tsvector(inDescription)));
  SELECT INTO lastInsertId currval('product_product_id_seq');

  INSERT INTO product_category (product_id, category_id) VALUES (lastInsertId,inCategoryId);

  UPDATE category SET category_quantity = category_quantity + 1 WHERE category_id = inCategoryId;

  RETURN lastInsertId;
END;
$_$;


ALTER FUNCTION public.product_add(bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]) OWNER TO postgres;

--
-- Name: product_update(bigint, bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION product_update(bigint, bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]) RETURNS bigint
    LANGUAGE plpgsql
    AS $_$
DECLARE
  inId              ALIAS FOR $1;
  inCategoryId      ALIAS FOR $2;
  inTitle           ALIAS FOR $3;
  inDescription     ALIAS FOR $4;
  inCode            ALIAS FOR $5;
  inPrice           ALIAS FOR $6;
  inPrice1          ALIAS FOR $7;
  inEnable          ALIAS FOR $8;
  inImgId           ALIAS FOR $9;

BEGIN

--IF (SELECT count(*) FROM img WHERE img_id = ALL (inImgId)) != array_length(inImgId,1)
--THEN RETURN -10;
--END IF;

  UPDATE product SET (
    product_title,
    product_description,
    product_code,
    product_price,
    product_price1,
    product_enable,
    product_img,
    search_vector)
  =
  (inTitle, inDescription,inCode,inPrice,inPrice1,inEnable,inImgId,
   (setweight(to_tsvector(inTitle), 'A') || to_tsvector(inDescription)))
  WHERE product_id = inId;

  UPDATE product_category SET category_id = inCategoryId WHERE product_id = inId;

  RETURN 1;
END;
$_$;


ALTER FUNCTION public.product_update(bigint, bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]) OWNER TO postgres;

--
-- Name: update_product_updated_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION update_product_updated_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
    NEW.product_updated = now();
    RETURN NEW;
  ELSE
    RETURN OLD;
  END IF;
END;
$$;


ALTER FUNCTION public.update_product_updated_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE account (
    account_id bigint NOT NULL,
    account_email character varying(255),
    account_nick character varying(255),
    account_password character varying(512) DEFAULT ''::character varying,
    account_phone character varying(64),
    account_provider character varying(512) DEFAULT ''::character varying NOT NULL,
    account_token character varying(512) DEFAULT ''::character varying NOT NULL,
    account_ban_reason character varying(255) DEFAULT ''::character varying NOT NULL,
    account_newpass character varying(64) DEFAULT ''::character varying NOT NULL,
    account_newpass_key character varying(64) DEFAULT ''::character varying NOT NULL,
    account_newpass_time timestamp without time zone DEFAULT now() NOT NULL,
    account_last_ip character varying(64) DEFAULT ''::character varying NOT NULL,
    account_last_logged timestamp without time zone DEFAULT now() NOT NULL,
    account_created timestamp without time zone DEFAULT now() NOT NULL,
    account_updated timestamp without time zone DEFAULT now() NOT NULL,
    account_position bigint,
    account_ban boolean DEFAULT false NOT NULL
);


ALTER TABLE account OWNER TO postgres;

--
-- Name: account_account_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE account_account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE account_account_id_seq OWNER TO postgres;

--
-- Name: account_account_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE account_account_id_seq OWNED BY account.account_id;


--
-- Name: activation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE activation (
    activation_id bigint NOT NULL,
    activation_email character varying(256) DEFAULT ''::character varying NOT NULL,
    activation_nick character varying(256) DEFAULT ''::character varying NOT NULL,
    activation_password character varying(128) DEFAULT ''::character varying NOT NULL,
    activation_key character varying(256) DEFAULT ''::character varying NOT NULL,
    activation_last_ip character varying(64) DEFAULT ''::character varying NOT NULL,
    activation_created bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL,
    activation_phone character varying(64) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE activation OWNER TO postgres;

--
-- Name: activation_activation_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE activation_activation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE activation_activation_id_seq OWNER TO postgres;

--
-- Name: activation_activation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE activation_activation_id_seq OWNED BY activation.activation_id;


--
-- Name: attempt; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE attempt (
    attempt_id bigint NOT NULL,
    attempt_ip character varying(64) DEFAULT ''::character varying NOT NULL,
    attempt_time bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL
);


ALTER TABLE attempt OWNER TO postgres;

--
-- Name: attempt_attempt_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE attempt_attempt_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE attempt_attempt_id_seq OWNER TO postgres;

--
-- Name: attempt_attempt_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE attempt_attempt_id_seq OWNED BY attempt.attempt_id;


--
-- Name: cart; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE cart (
    cart_id character(128) DEFAULT ''::bpchar NOT NULL,
    cart_product bigint NOT NULL,
    cart_quantity bigint DEFAULT 0 NOT NULL,
    cart_buy_now boolean DEFAULT true NOT NULL,
    cart_created timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cart OWNER TO postgres;

--
-- Name: category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE category (
    category_id bigint NOT NULL,
    category_sort bigint DEFAULT 100 NOT NULL,
    category_title character varying(255) DEFAULT ''::character varying NOT NULL,
    category_description text DEFAULT ''::text NOT NULL,
    category_enable boolean DEFAULT true NOT NULL,
    category_img character varying(255)[] DEFAULT (ARRAY[]::character varying[])::character varying(128)[] NOT NULL,
    category_quantity bigint DEFAULT 0 NOT NULL,
    category_parent bigint,
    category_lang character varying(12) DEFAULT 'en'::character varying NOT NULL
);


ALTER TABLE category OWNER TO postgres;

--
-- Name: category_category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE category_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE category_category_id_seq OWNER TO postgres;

--
-- Name: category_category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE category_category_id_seq OWNED BY category.category_id;


--
-- Name: coupon; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE coupon (
    coupon_id bigint NOT NULL,
    coupon_code character varying(512) DEFAULT ''::character varying NOT NULL,
    coupon_start_time timestamp without time zone,
    coupon_end_time timestamp without time zone,
    coupon_used bigint DEFAULT 0 NOT NULL,
    coupon_used_limit bigint DEFAULT 1 NOT NULL,
    coupon_per_product bigint DEFAULT 1 NOT NULL
);


ALTER TABLE coupon OWNER TO postgres;

--
-- Name: coupon_coupon_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE coupon_coupon_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE coupon_coupon_id_seq OWNER TO postgres;

--
-- Name: coupon_coupon_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE coupon_coupon_id_seq OWNED BY coupon.coupon_id;


--
-- Name: img; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE img (
    img_id bigint NOT NULL,
    img_name character varying(255) DEFAULT ''::character varying NOT NULL,
    img_path character varying(255) DEFAULT ''::character varying NOT NULL,
    img_url character varying(255) DEFAULT ''::character varying NOT NULL,
    img_thumb character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE img OWNER TO postgres;

--
-- Name: img_img_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE img_img_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE img_img_id_seq OWNER TO postgres;

--
-- Name: img_img_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE img_img_id_seq OWNED BY img.img_id;


--
-- Name: permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE permission (
    permission_id bigint NOT NULL,
    permission_data text DEFAULT ''::text NOT NULL,
    permission_position bigint
);


ALTER TABLE permission OWNER TO postgres;

--
-- Name: permission_permission_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE permission_permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE permission_permission_id_seq OWNER TO postgres;

--
-- Name: permission_permission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE permission_permission_id_seq OWNED BY permission.permission_id;


--
-- Name: position; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE "position" (
    position_id bigint NOT NULL,
    position_title character varying(256) DEFAULT ''::character varying NOT NULL,
    position_parent bigint,
    position_sort integer DEFAULT 100 NOT NULL,
    position_enable boolean DEFAULT true NOT NULL
);


ALTER TABLE "position" OWNER TO postgres;

--
-- Name: position_position_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE position_position_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE position_position_id_seq OWNER TO postgres;

--
-- Name: position_position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE position_position_id_seq OWNED BY "position".position_id;


--
-- Name: product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE product (
    product_id bigint NOT NULL,
    product_price numeric(20,2) DEFAULT 0.00 NOT NULL,
    product_title character varying(255) DEFAULT ''::character varying NOT NULL,
    product_description character varying(1024) DEFAULT ''::character varying NOT NULL,
    product_created timestamp without time zone DEFAULT now() NOT NULL,
    product_updated timestamp without time zone DEFAULT now() NOT NULL,
    product_enable boolean DEFAULT true NOT NULL,
    search_vector tsvector,
    product_img character varying(255)[] DEFAULT (ARRAY[]::character varying[])::character varying(255)[] NOT NULL,
    product_code character varying(255) DEFAULT ''::character varying NOT NULL,
    product_price1 numeric(20,2) DEFAULT 0.00 NOT NULL,
    product_watch bigint DEFAULT 0 NOT NULL,
    product_like bigint[] DEFAULT ARRAY[]::bigint[] NOT NULL,
    product_soled bigint DEFAULT 0 NOT NULL
);


ALTER TABLE product OWNER TO postgres;

--
-- Name: product_category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE product_category (
    product_id bigint NOT NULL,
    category_id bigint NOT NULL
);


ALTER TABLE product_category OWNER TO postgres;

--
-- Name: product_product_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE product_product_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE product_product_id_seq OWNER TO postgres;

--
-- Name: product_product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE product_product_id_seq OWNED BY product.product_id;


--
-- Name: session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE session (
    session_id character varying(512) NOT NULL,
    session_ip character varying(64) DEFAULT ''::character varying NOT NULL,
    session_timestamp bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL,
    session_account bigint,
    session_email character varying(255),
    session_nick character varying(255),
    session_data text DEFAULT ''::text NOT NULL,
    session_agent character varying(255) DEFAULT ''::character varying NOT NULL,
    session_data1 text DEFAULT ''::text NOT NULL,
    session_data2 text DEFAULT ''::text NOT NULL,
    session_phone character varying(64),
    session_provider character varying(255) DEFAULT ''::character varying,
    session_token character varying(255) DEFAULT ''::character varying,
    session_position bigint
);


ALTER TABLE session OWNER TO postgres;

--
-- Name: shipment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE shipment (
    shipment_id bigint NOT NULL,
    shipment_time timestamp without time zone DEFAULT now() NOT NULL,
    shipment_customer_name character varying(128) DEFAULT ''::character varying NOT NULL,
    shipment_email character varying(255) DEFAULT ''::character varying NOT NULL,
    shipment_country character varying(64) DEFAULT ''::character varying NOT NULL,
    shipment_zip character varying(64) DEFAULT ''::character varying NOT NULL,
    shipment_state character varying(64) DEFAULT ''::character varying NOT NULL,
    shipment_city character varying(64) DEFAULT ''::character varying NOT NULL,
    shipment_address character varying(128) DEFAULT ''::character varying NOT NULL,
    shipment_phone character varying(64) DEFAULT ''::character varying NOT NULL,
    shipment_status integer DEFAULT 0 NOT NULL,
    shipment_comment character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE shipment OWNER TO postgres;

--
-- Name: shipment_shipment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE shipment_shipment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE shipment_shipment_id_seq OWNER TO postgres;

--
-- Name: shipment_shipment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE shipment_shipment_id_seq OWNED BY shipment.shipment_id;


--
-- Name: account account_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account ALTER COLUMN account_id SET DEFAULT nextval('account_account_id_seq'::regclass);


--
-- Name: activation activation_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY activation ALTER COLUMN activation_id SET DEFAULT nextval('activation_activation_id_seq'::regclass);


--
-- Name: attempt attempt_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY attempt ALTER COLUMN attempt_id SET DEFAULT nextval('attempt_attempt_id_seq'::regclass);


--
-- Name: category category_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY category ALTER COLUMN category_id SET DEFAULT nextval('category_category_id_seq'::regclass);


--
-- Name: coupon coupon_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY coupon ALTER COLUMN coupon_id SET DEFAULT nextval('coupon_coupon_id_seq'::regclass);


--
-- Name: img img_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY img ALTER COLUMN img_id SET DEFAULT nextval('img_img_id_seq'::regclass);


--
-- Name: permission permission_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission ALTER COLUMN permission_id SET DEFAULT nextval('permission_permission_id_seq'::regclass);


--
-- Name: position position_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "position" ALTER COLUMN position_id SET DEFAULT nextval('position_position_id_seq'::regclass);


--
-- Name: product product_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product ALTER COLUMN product_id SET DEFAULT nextval('product_product_id_seq'::regclass);


--
-- Name: shipment shipment_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY shipment ALTER COLUMN shipment_id SET DEFAULT nextval('shipment_shipment_id_seq'::regclass);


--
-- Data for Name: account; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO account VALUES (38, 'ilyaran1@mail.ru', NULL, 'fc7d504552eec1104873a1d86c12f360145e59153cfd678e57d252857def34f5', NULL, '', '', '', '', '', '2017-05-27 09:08:46.909604', '::1', '2017-05-29 13:38:39.761756', '2017-05-27 09:08:46.909604', '2017-05-27 09:08:46.909604', 2, false);
INSERT INTO account VALUES (44, 'ilyaraaan@mail.ru', 'player777www', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '+77026440799', '', '', '', '', '', '2017-05-28 16:13:29.396083', '', '2017-05-28 16:13:29.396083', '2017-05-28 16:13:29.396083', '2017-05-28 16:13:29.396083', NULL, false);
INSERT INTO account VALUES (43, NULL, 'player6', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', NULL, '', '', '', '', '', '2017-05-28 16:10:10.613654', '', '2017-05-28 16:10:10.613654', '2017-05-28 16:10:10.613654', '2017-05-28 16:10:10.613654', NULL, false);
INSERT INTO account VALUES (42, NULL, 'player5', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', NULL, '', '', '', '', '', '2017-05-28 16:09:00.927459', '', '2017-05-28 16:09:00.927459', '2017-05-28 16:09:00.927459', '2017-05-28 16:09:00.927459', NULL, true);
INSERT INTO account VALUES (41, NULL, 'player4', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', NULL, '', '', '', '', '', '2017-05-28 16:06:33.81613', '', '2017-05-28 16:06:33.81613', '2017-05-28 16:06:33.81613', '2017-05-28 16:06:33.81613', NULL, true);
INSERT INTO account VALUES (40, 'ilyaran3@mail.ru', NULL, 'a737481c2a91b0dc3db235491c85c4a997fb587c02ce04391220581db04fe64c', NULL, '', '', '', '', '', '2017-05-28 11:44:15.946587', '', '2017-05-28 11:44:15.946587', '2017-05-28 11:44:15.946587', '2017-05-28 11:44:15.946587', NULL, false);
INSERT INTO account VALUES (39, 'ilyaran2@mail.ru', NULL, 'd771e6d5feec19bf37a04232285933d9dcf202757a8f4af53771b348ead309a6', NULL, '', '', '', '', '', '2017-05-27 09:31:26.610499', '::1', '2017-05-27 17:26:20.286622', '2017-05-27 09:31:26.610499', '2017-05-27 09:31:26.610499', NULL, false);
INSERT INTO account VALUES (1, 'ilyaran@mail.ru', 'ilyaran', 'f3635736ad71032133dfbbf8dcdd136edf5b938d0c19b2bf2cd0f44747b0678e', '+77058436633', '', '', '', '', '', '2017-05-26 20:38:46.129352', '::1', '2017-05-29 15:07:07.406515', '2017-05-26 20:38:46.129352', '2017-05-26 20:38:46.129352', 3, false);


--
-- Name: account_account_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('account_account_id_seq', 44, true);


--
-- Data for Name: activation; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: activation_activation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('activation_activation_id_seq', 4, true);


--
-- Data for Name: attempt; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: attempt_attempt_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('attempt_attempt_id_seq', 1, true);


--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 29, 2, true, '2017-04-28 15:42:42.603732');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 27, 1, true, '2017-04-28 15:57:01.561923');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 25, 1, true, '2017-04-28 15:59:30.157029');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 24, 1, true, '2017-04-28 15:59:49.330372');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 26, 4, true, '2017-04-28 15:59:21.41243');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 23, 1, true, '2017-04-28 16:25:43.883635');
INSERT INTO cart VALUES ('09c85b08d92d5d1fcd2c544fa4ece1511c84a025e7d2613826bb7b1d8daba316                                                                ', 28, 1, true, '2017-04-28 16:30:14.216816');
INSERT INTO cart VALUES ('6b5eae5f37f54ad36e07dd63c36843d5e516d7caa8adf0fcd8adf7eddeebf6ba                                                                ', 28, 1, true, '2017-04-30 11:40:54.518462');
INSERT INTO cart VALUES ('6b5eae5f37f54ad36e07dd63c36843d5e516d7caa8adf0fcd8adf7eddeebf6ba                                                                ', 30, 1, true, '2017-04-30 11:40:54.518462');
INSERT INTO cart VALUES ('6b5eae5f37f54ad36e07dd63c36843d5e516d7caa8adf0fcd8adf7eddeebf6ba                                                                ', 27, 5, true, '2017-04-30 11:40:54.518462');
INSERT INTO cart VALUES ('6b5eae5f37f54ad36e07dd63c36843d5e516d7caa8adf0fcd8adf7eddeebf6ba                                                                ', 29, 1, true, '2017-04-30 11:40:54.518462');
INSERT INTO cart VALUES ('dee0673b12408c0d95501ed9a3c80a2f95ec7dc46fa306d4052aa3202a303584                                                                ', 33, 1, true, '2017-05-29 11:24:26.29628');
INSERT INTO cart VALUES ('898e181cfb20842a11000bf593ba9a841a5bf714a2f60e339d4ba850c8fb562a                                                                ', 26, 1, true, '2017-04-29 16:14:53.86801');
INSERT INTO cart VALUES ('898e181cfb20842a11000bf593ba9a841a5bf714a2f60e339d4ba850c8fb562a                                                                ', 25, 7, true, '2017-04-29 16:14:53.86801');
INSERT INTO cart VALUES ('dee0673b12408c0d95501ed9a3c80a2f95ec7dc46fa306d4052aa3202a303584                                                                ', 32, 1, true, '2017-05-29 11:59:39.767803');
INSERT INTO cart VALUES ('dee0673b12408c0d95501ed9a3c80a2f95ec7dc46fa306d4052aa3202a303584                                                                ', 36, 1, true, '2017-05-29 11:59:45.78594');
INSERT INTO cart VALUES ('dee0673b12408c0d95501ed9a3c80a2f95ec7dc46fa306d4052aa3202a303584                                                                ', 29, 5, true, '2017-05-29 11:24:41.759784');


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO category VALUES (13, 140, 'Food', '', true, '{/assets/uploads/roxy-fileman-logo.gif}', 3, NULL, '');
INSERT INTO category VALUES (9, 9225, 'Computers', '<p><img alt="" src="http://localhost:3001/assets/uploads/1e43d48s-960.jpg" style="height:533px; width:800px" /></p>
', true, '{/assets/uploads/1e43d48s-960.jpg}', 0, NULL, '');
INSERT INTO category VALUES (20, 100, 'Notebooks', '<p><em><strong>Washing machines</strong></em></p>

<p><em><strong><img alt="" src="http://localhost:3001/assets/uploads/DSC_2987.jpg" style="height:469px; width:700px" /></strong></em></p>
', true, '{/assets/uploads/1e43d48s-960.jpg,/assets/uploads/21.gif}', 2, 9, '');
INSERT INTO category VALUES (22, 100, 'Ultrabooks', '', true, '{/assets/uploads/bag.jpg,/assets/uploads/bo.jpg}', 0, 21, '');
INSERT INTO category VALUES (25, 110, 'Ultra Gadgets', '<p><img alt="" src="http://localhost:3001/assets/uploads/pic11.jpg" style="height:285px; width:210px" /></p>
', true, '{/assets/uploads/images555.jpg,/assets/uploads/images.jpg}', 0, 22, '');
INSERT INTO category VALUES (10, 18025, 'Laptops', '<p><img alt="" src="http://localhost:3001/assets/uploads/roxy-fileman-logo.gif" style="height:127px; width:289px" /></p>
', true, '{/assets/uploads/pic5.jpg,/assets/uploads/pic7.jpg}', 0, NULL, '');
INSERT INTO category VALUES (11, 150, 'Furniture', '', true, '{/assets/uploads/1e43d48s-960.jpg}', 0, NULL, '');
INSERT INTO category VALUES (23, 500, 'Animals pets', '', true, '{/assets/uploads/Documents/uugugugu.jpg,/assets/uploads/Documents/000006801x.jpg}', 6, NULL, '');
INSERT INTO category VALUES (21, 110, 'Tablets ddd', '', true, '{}', 1, 20, '');
INSERT INTO category VALUES (24, 101, 'Transport', '', true, '{/assets/uploads/images.jpg,"/assets/uploads/скачанные файлы.jpg",/assets/uploads/21.gif}', 5, NULL, '');


--
-- Name: category_category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('category_category_id_seq', 28, true);


--
-- Data for Name: coupon; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: coupon_coupon_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('coupon_coupon_id_seq', 1, false);


--
-- Data for Name: img; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO img VALUES (1, 'boss', '', '', '');
INSERT INTO img VALUES (2, 'boss', '', '', '');


--
-- Name: img_img_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('img_img_id_seq', 2, true);


--
-- Data for Name: permission; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO permission VALUES (8, 'admin', 3);
INSERT INTO permission VALUES (9, 'admin', 4);
INSERT INTO permission VALUES (10, 'admin ddd', 2);
INSERT INTO permission VALUES (11, 'category fdff', 8);
INSERT INTO permission VALUES (12, 'permission', 8);


--
-- Name: permission_permission_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('permission_permission_id_seq', 13, true);


--
-- Data for Name: position; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "position" VALUES (3, 'boss', NULL, 100, true);
INSERT INTO "position" VALUES (4, 'manager', 3, 100, true);
INSERT INTO "position" VALUES (2, 'admin  ', 4, 95, true);
INSERT INTO "position" VALUES (1, 'user', 2, 100, true);
INSERT INTO "position" VALUES (6, 'sub admin lll', 2, 110, true);
INSERT INTO "position" VALUES (8, 'under sub admin', 6, 100, true);
INSERT INTO "position" VALUES (10, '333 under sub', 6, 120, true);
INSERT INTO "position" VALUES (7, 'sub admin iooii', 2, 120, true);
INSERT INTO "position" VALUES (15, 'z3xx dfsd dfxdf', 4, 100, true);
INSERT INTO "position" VALUES (12, 'zxx dfsd dfxdf', 4, 100, true);
INSERT INTO "position" VALUES (13, 'z1xx dfsd dfxdf', 4, 100, true);
INSERT INTO "position" VALUES (14, 'z2xx dfsd dfxdf', 4, 100, true);


--
-- Name: position_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('position_position_id_seq', 18, true);


--
-- Data for Name: product; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO product VALUES (21, 102.00, '2558 My Product', '', '2017-04-25 17:48:37.961616', '2017-05-06 13:13:16.332519', true, '''2558'':1A ''product'':3A', '{/assets/uploads/bag.jpg,/assets/uploads/1e43d48s-960.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (26, 287.00, '765 Some pet', '', '2017-04-26 07:53:04.817909', '2017-05-01 15:06:24.220413', true, '''765'':1A ''pet'':3A', '{/assets/uploads/pic.jpg,/assets/uploads/pic10.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (23, 235.00, 'Some pet good', '<h3><var><em><strong>Возвращаясь к пункту про структуру директории, gb трактует всё, что находится в&nbsp;src/, как код вашего </strong></em></var></h3>

<hr />
<h3><var><em><strong>проекта. Все зависимые пакаджи устанавливаются в директор</strong></em></var></h3>

<hr />
<h3><var><em><strong>ию&nbsp;vendor/&nbsp;и именно оттуда код берется при сборке с помощью gb.</strong></em></var></h3>
', '2017-04-26 07:52:43.79536', '2017-05-01 15:05:57.121359', true, '''gb'':10,38 ''good'':3A ''pet'':2A ''src'':16 ''vendor'':28 ''берет'':33 ''ваш'':19 ''возвра'':4 ''всё'':12 ''директор'':9,26 ''зависим'':22 ''и'':27 ''имен'':30 ''код'':18,32 ''наход'':14 ''оттуд'':31 ''пакадж'':23 ''помощ'':37 ''проект'':20 ''пункт'':6 ''сборк'':35 ''структур'':8 ''тракт'':11 ''устанавлива'':24', '{/assets/uploads/ch.jpg,/assets/uploads/pi.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (25, 25684.18, '554 Some pet dddd', '', '2017-04-26 07:52:57.60257', '2017-04-29 17:55:59.669302', true, '''554'':1A ''dddd'':4A ''pet'':3A', '{/assets/uploads/pic11.jpg,/assets/uploads/pic2.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (36, 6565.36, '11 Watch Some vvv', '', '2017-05-26 12:54:55.44247', '2017-05-26 12:59:30.508373', true, '''11'':1A ''vvv'':4A ''watch'':2A', '{}', '', 6445.78, 0, '{}', 0);
INSERT INTO product VALUES (32, 0.00, 'some product gg', '', '2017-05-08 11:14:08.67008', '2017-05-26 12:43:48.156194', true, '''product'':2A', '{/assets/uploads/1e43d48s-960.jpg,/assets/uploads/21.gif}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (33, 6565.36, 'Watch Some vvv', '', '2017-05-26 12:41:26.912277', '2017-05-26 12:43:08.875805', true, '''vvv'':3A ''watch'':1A', '{/assets/uploads/pic11.jpg,/assets/uploads/images.jpg}', '', 6445.78, 0, '{}', 0);
INSERT INTO product VALUES (28, 378.00, 'Parfume ccc', '', '2017-04-27 15:34:09.240468', '2017-05-01 15:07:17.379328', true, '''ccc'':2A ''parfum'':1A', '{/assets/uploads/pic4.jpg,/assets/uploads/s4.jpg,/assets/uploads/sh.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (27, 339.00, 'Shoes nike', '', '2017-04-27 15:33:40.925847', '2017-05-01 15:06:59.854226', true, '''nike'':2A ''shoe'':1A', '{/assets/uploads/sh.jpg,/assets/uploads/s4.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (30, 312.00, 'Girls ccc', '', '2017-04-27 15:35:17.124572', '2017-05-01 15:06:40.654448', true, '''ccc'':2A ''girl'':1A', '{/assets/uploads/pic12.jpg,/assets/uploads/pi4.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (24, 186.00, '344 Some pet', '', '2017-04-26 07:52:50.973054', '2017-05-01 15:05:39.843489', true, '''344'':1A ''pet'':3A', '{/assets/uploads/roxy-fileman-logo.gif,/assets/uploads/000006801x.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (22, 152.00, 'My ProductOOOKKK', '', '2017-04-26 07:43:57.314561', '2017-05-01 15:05:23.20969', true, '''productoookkk'':2A', '{/assets/uploads/21.gif,/assets/uploads/11395_original.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (29, 647.21, 'Bag ccc', '', '2017-04-27 15:34:40.735855', '2017-05-01 10:43:54.115379', true, '''bag'':1A ''ccc'':2A', '{/assets/uploads/s4.jpg,/assets/uploads/sh.jpg,/assets/uploads/pic4.jpg}', '', 589.36, 0, '{}', 0);
INSERT INTO product VALUES (31, 568.23, 'Washing machine', '<p><img alt="" src="http://localhost:3001/assets/uploads/21.gif" style="height:272px; width:615px" />&nbsp;</p>

<h1>К ЧЕМУ СНИТСЯ, ЧТО УМЕР ЖИВОЙ. ТОЛКОВАНИЕ СНОВ</h1>

<blockquote>
<h2 style="font-style:italic;">Сон, в котором вы видите, как умер кто-то из ваших живых знакомых или родственников, означает, прежде всего, долголетие этого человека. Однако его правильное толкование зависит от того, кто именно умер. Если умерли ваши знакомые, сон предупреждает о возможности серьёзного конфликта с ними из-за вашей несдержанности или легкомыслия. Вам следует контролировать своё поведение и слова. Если умер друг, по какой-то причине отношения с ним могут прерваться.</h2>
</blockquote>
', '2017-05-01 15:37:00.218383', '2017-05-01 15:37:00.218383', true, '''machin'':2A ''wash'':1A ''ваш'':22,45,58 ''вид'':15 ''возможн'':50 ''долголет'':30 ''друг'':71 ''жив'':8,23 ''завис'':37 ''знаком'':24,46 ''из-з'':55 ''имен'':41 ''какой-т'':73 ''контролирова'':64 ''конфликт'':52 ''котор'':13 ''кто-т'':18 ''легкомысл'':61 ''могут'':80 ''несдержан'':59 ''ним'':54 ''однак'':33 ''означа'':27 ''отношен'':77 ''поведен'':66 ''правильн'':35 ''предупрежда'':48 ''прежд'':28 ''прерва'':81 ''причин'':76 ''родственник'':26 ''своё'':65 ''серьёзн'':51 ''след'':63 ''слов'':68 ''снит'':5 ''снов'':10 ''сон'':11,47 ''толкован'':9,36 ''умер'':7,17,42,70 ''умерл'':44 ''человек'':32 ''чем'':4', '{/assets/uploads/1e43d48s-960.jpg,/assets/uploads/000006801x.jpg,/assets/uploads/wat.jpg,/assets/uploads/11395_original.jpg}', '', 0.00, 0, '{}', 0);


--
-- Data for Name: product_category; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO product_category VALUES (25, 23);
INSERT INTO product_category VALUES (29, 20);
INSERT INTO product_category VALUES (22, 13);
INSERT INTO product_category VALUES (24, 23);
INSERT INTO product_category VALUES (23, 23);
INSERT INTO product_category VALUES (26, 23);
INSERT INTO product_category VALUES (30, 20);
INSERT INTO product_category VALUES (27, 23);
INSERT INTO product_category VALUES (28, 23);
INSERT INTO product_category VALUES (31, 21);
INSERT INTO product_category VALUES (21, 13);
INSERT INTO product_category VALUES (32, 24);
INSERT INTO product_category VALUES (33, 24);
INSERT INTO product_category VALUES (36, 24);


--
-- Name: product_product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('product_product_id_seq', 36, true);


--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO session VALUES ('6e43effc60d3dd5dd6885c5ea5e68df6b677ef3c0023789014b55717a2f17e72', '::1', 1496072114, NULL, NULL, NULL, 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', '', '', NULL, '', '', NULL);


--
-- Data for Name: shipment; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: shipment_shipment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('shipment_shipment_id_seq', 1, false);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_pkey PRIMARY KEY (account_id);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY category
    ADD CONSTRAINT category_pkey PRIMARY KEY (category_id);


--
-- Name: coupon coupon_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY coupon
    ADD CONSTRAINT coupon_pkey PRIMARY KEY (coupon_id);


--
-- Name: img img_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY img
    ADD CONSTRAINT img_pkey PRIMARY KEY (img_id);


--
-- Name: permission permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission
    ADD CONSTRAINT permission_pkey PRIMARY KEY (permission_id);


--
-- Name: product_category pk_product_id_category_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product_category
    ADD CONSTRAINT pk_product_id_category_id PRIMARY KEY (product_id, category_id);


--
-- Name: position position_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "position"
    ADD CONSTRAINT position_pkey PRIMARY KEY (position_id);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product
    ADD CONSTRAINT product_pkey PRIMARY KEY (product_id);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_pkey PRIMARY KEY (session_id);


--
-- Name: shipment shipment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY shipment
    ADD CONSTRAINT shipment_pkey PRIMARY KEY (shipment_id);


--
-- Name: account_account_email_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_email_uindex ON account USING btree (account_email);


--
-- Name: account_account_nick_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_nick_uindex ON account USING btree (account_nick);


--
-- Name: account_account_phone_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_phone_uindex ON account USING btree (account_phone);


--
-- Name: product update_product_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_product_updated BEFORE UPDATE ON product FOR EACH ROW EXECUTE PROCEDURE update_product_updated_column();


--
-- Name: account account_position_position_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_position_position_id_fk FOREIGN KEY (account_position) REFERENCES "position"(position_id) ON DELETE SET NULL;


--
-- Name: category category_category_category_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY category
    ADD CONSTRAINT category_category_category_id_fk FOREIGN KEY (category_parent) REFERENCES category(category_id) ON DELETE SET NULL;


--
-- Name: cart fk_cart_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY cart
    ADD CONSTRAINT fk_cart_product FOREIGN KEY (cart_product) REFERENCES product(product_id) ON DELETE CASCADE;


--
-- Name: product_category fk_category_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product_category
    ADD CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(category_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: permission fk_permission_position_position_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission
    ADD CONSTRAINT fk_permission_position_position_id FOREIGN KEY (permission_position) REFERENCES "position"(position_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: product_category fk_product_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product_category
    ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES product(product_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: position position_position_position_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "position"
    ADD CONSTRAINT position_position_position_id_fk FOREIGN KEY (position_parent) REFERENCES "position"(position_id) ON DELETE SET NULL;


--
-- Name: session session_account_account_email_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_email_fk FOREIGN KEY (session_email) REFERENCES account(account_email) ON DELETE SET NULL;


--
-- Name: session session_account_account_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_id_fk FOREIGN KEY (session_account) REFERENCES account(account_id) ON DELETE SET NULL;


--
-- Name: session session_account_account_nick_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_nick_fk FOREIGN KEY (session_nick) REFERENCES account(account_nick) ON DELETE SET NULL;


--
-- Name: session session_account_account_phone_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_phone_fk FOREIGN KEY (session_phone) REFERENCES account(account_phone) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_position_position_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_position_position_id_fk FOREIGN KEY (session_position) REFERENCES "position"(position_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

