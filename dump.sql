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
-- Name: delete_expired_activation(bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION delete_expired_activation(bigint) RETURNS void
    LANGUAGE plpgsql
    AS $_$

DECLARE activationLifeSpan ALIAS FOR $1;

BEGIN

  DELETE FROM user_temp

  WHERE date_part('epoch',CURRENT_TIMESTAMP)::bigint - created > activationLifeSpan;

END;

$_$;


ALTER FUNCTION public.delete_expired_activation(bigint) OWNER TO postgres;

--
-- Name: delete_expired_sessions(bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION delete_expired_sessions(bigint) RETURNS void
    LANGUAGE plpgsql
    AS $_$
DECLARE sessionLifeSpan ALIAS FOR $1;
BEGIN
  DELETE FROM sessions
  WHERE date_part('epoch',CURRENT_TIMESTAMP)::bigint - timestamp > sessionLifeSpan;
END;
$_$;


ALTER FUNCTION public.delete_expired_sessions(bigint) OWNER TO postgres;

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
    account_password character varying(512) DEFAULT ''::character varying NOT NULL,
    account_phone character varying(64) DEFAULT ''::character varying NOT NULL,
    account_fist_name character varying(255) DEFAULT ''::character varying NOT NULL,
    account_last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    account_img character varying(255) DEFAULT NULL::character varying,
    account_provider character varying(512) DEFAULT ''::character varying NOT NULL,
    account_token character varying(512) DEFAULT ''::character varying NOT NULL,
    account_banned smallint DEFAULT 0 NOT NULL,
    account_ban_reason character varying(255) DEFAULT ''::character varying NOT NULL,
    account_newpass character varying(64) DEFAULT ''::character varying NOT NULL,
    account_newpass_key character varying(64) DEFAULT ''::character varying NOT NULL,
    account_newpass_time timestamp without time zone DEFAULT now() NOT NULL,
    account_last_ip character varying(512) DEFAULT ''::character varying NOT NULL,
    account_last_logged timestamp without time zone DEFAULT now() NOT NULL,
    account_created timestamp without time zone DEFAULT now() NOT NULL,
    account_birth date DEFAULT now() NOT NULL,
    account_state character varying(64) DEFAULT ''::character varying NOT NULL,
    account_city character varying(64) DEFAULT ''::character varying NOT NULL,
    account_skype character varying(255) DEFAULT ''::character varying NOT NULL,
    account_steam_id character varying(255) DEFAULT ''::character varying NOT NULL,
    account_balance bigint DEFAULT (0)::bigint NOT NULL,
    account_updated timestamp without time zone DEFAULT now() NOT NULL,
    account_position bigint DEFAULT 1 NOT NULL
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
    test timestamp without time zone DEFAULT now() NOT NULL
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
-- Name: login_attempts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE login_attempts (
    attempt_id bigint NOT NULL,
    ip_address character varying(64) DEFAULT ''::character varying NOT NULL,
    attempt_time bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL
);


ALTER TABLE login_attempts OWNER TO postgres;

--
-- Name: login_attempts_attempt_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE login_attempts_attempt_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE login_attempts_attempt_id_seq OWNER TO postgres;

--
-- Name: login_attempts_attempt_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE login_attempts_attempt_id_seq OWNED BY login_attempts.attempt_id;


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
    position_parent bigint DEFAULT (0)::bigint NOT NULL,
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
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE sessions (
    id character varying(512) NOT NULL,
    ip_address character varying(64) DEFAULT ''::character varying NOT NULL,
    "timestamp" bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL,
    account_id bigint DEFAULT 0 NOT NULL,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    nick character varying(255) DEFAULT ''::character varying NOT NULL,
    data text DEFAULT ''::text NOT NULL,
    user_agent character varying(255) DEFAULT ''::character varying NOT NULL,
    is_flash boolean DEFAULT false NOT NULL,
    position_id bigint DEFAULT 0 NOT NULL,
    position_title character varying(255) DEFAULT ''::character varying NOT NULL,
    permission_data character varying(512) DEFAULT ''::character varying NOT NULL,
    balance bigint DEFAULT 0 NOT NULL,
    phone character varying(64) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE sessions OWNER TO postgres;

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
-- Name: user_temp; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE user_temp (
    id bigint NOT NULL,
    user_email character varying(256) DEFAULT ''::character varying NOT NULL,
    user_name character varying(256) DEFAULT ''::character varying NOT NULL,
    user_password character varying(128) DEFAULT ''::character varying NOT NULL,
    activation_key character varying(256) DEFAULT ''::character varying NOT NULL,
    last_ip character varying(64) DEFAULT ''::character varying NOT NULL,
    created bigint DEFAULT (date_part('epoch'::text, now()))::bigint NOT NULL,
    user_phone character varying(64) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE user_temp OWNER TO postgres;

--
-- Name: user_temp_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_temp_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_temp_id_seq OWNER TO postgres;

--
-- Name: user_temp_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_temp_id_seq OWNED BY user_temp.id;


--
-- Name: account account_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account ALTER COLUMN account_id SET DEFAULT nextval('account_account_id_seq'::regclass);


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
-- Name: login_attempts attempt_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY login_attempts ALTER COLUMN attempt_id SET DEFAULT nextval('login_attempts_attempt_id_seq'::regclass);


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
-- Name: user_temp id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_temp ALTER COLUMN id SET DEFAULT nextval('user_temp_id_seq'::regclass);


--
-- Data for Name: account; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO account VALUES (8, 'il.aranov@gmail.com', '', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-03 00:00:00', '::1', '2017-04-04 15:13:45.053954', '2017-04-03 00:00:00', '2017-04-03', '', '', '', '', 0, '2017-04-10 15:40:37.841358', 2);
INSERT INTO account VALUES (14, '', 'player1', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-07 00:00:00', '', '2017-04-07 15:56:45.173638', '2017-04-07 15:56:45.173638', '2017-04-07', '', '', '', '', 0, '2017-04-07 16:46:02.92524', 1);
INSERT INTO account VALUES (13, '', 'newplayer12445', '7043e022c3726c01d1d5ebfa5c17d6f63cfee3aa7c6f1dfb31c119596979d0a5', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-06 00:00:00', '::1', '2017-04-07 16:51:51.262466', '2017-04-06 12:27:56.028942', '2017-04-06', '', '', '', '', 0, '2017-04-07 16:18:43.706946', 1);
INSERT INTO account VALUES (32, 'ilyadddran@mail.ru', '', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:38:28.145921', '2017-04-12 13:38:28.145921', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:38:28.145921', 2);
INSERT INTO account VALUES (20, '', 'player66', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '', '2017-04-08 13:27:09.470382', '2017-04-08 13:27:09.470382', '2017-04-08', '', '', '', '', 0, '2017-04-08 13:27:09.470382', 1);
INSERT INTO account VALUES (1, 'ilyaran@mail.ru', '', 'ae57c2cac1b184d0c639678b1386ee91af322297ef094dd77125d86e4daef6c4', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-02 00:00:00', '', '2017-04-02 00:00:00', '2017-04-02 00:00:00', '2017-04-02', '', '', '', '', 0, '2017-04-06 17:23:49.905039', 2);
INSERT INTO account VALUES (36, 'i-arajhn@mail.ru', 'player987', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:53:14.389987', '2017-04-12 13:53:14.389987', '2017-04-12', '', '', '', '', 0, '2017-04-12 14:52:44.41347', 3);
INSERT INTO account VALUES (27, '', 'player777', '', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:11:07.871447', '2017-04-12 13:11:07.871447', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:11:07.871447', 2);
INSERT INTO account VALUES (33, '', 'player111', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:38:50.759509', '2017-04-12 13:38:50.759509', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:38:50.759509', 2);
INSERT INTO account VALUES (28, '', 'player555', '', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:11:55.996985', '2017-04-12 13:11:55.996985', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:11:55.996985', 2);
INSERT INTO account VALUES (34, 'ilyyyyaran@mail.ru', 'player101', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:39:10.015661', '2017-04-12 13:39:10.015661', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:39:10.015661', 2);
INSERT INTO account VALUES (26, '', 'player45', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-10 00:00:00', '', '2017-04-10 16:09:38.07578', '2017-04-10 16:09:38.07578', '2017-04-10', '', '', '', '', 0, '2017-04-10 16:10:19.213573', 1);
INSERT INTO account VALUES (23, 'jfjsdfjdf@jfjjfjf.rt', 'player62', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '', '2017-04-08 13:30:32.000273', '2017-04-08 13:30:32.000273', '2017-04-08', '', '', '', '', 0, '2017-04-08 13:30:32.000273', 1);
INSERT INTO account VALUES (29, '', 'player444', '', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:20:00.210127', '2017-04-12 13:20:00.210127', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:20:00.210127', 2);
INSERT INTO account VALUES (18, 'i-aran@mail.ru', 'player6', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '::1', '2017-04-08 09:48:22.236521', '2017-04-08 09:19:46.515129', '2017-04-08', '', '', '', '', 0, '2017-04-08 09:19:46.515129', 1);
INSERT INTO account VALUES (35, '', 'player258', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:45:51.47693', '2017-04-12 13:45:51.47693', '2017-04-12', '', '', '', '', 0, '2017-04-12 14:53:45.864807', 3);
INSERT INTO account VALUES (21, '', 'player64', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '', '2017-04-08 13:27:35.278542', '2017-04-08 13:27:35.278542', '2017-04-08', '', '', '', '', 0, '2017-04-10 15:33:53.497937', 2);
INSERT INTO account VALUES (30, '', 'player333', '', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:21:29.805553', '2017-04-12 13:21:29.805553', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:21:29.805553', 2);
INSERT INTO account VALUES (24, 'aran@mail.ru', 'player44', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '', '2017-04-08 13:34:24.387065', '2017-04-08 13:34:24.387065', '2017-04-08', '', '', '', '', 0, '2017-04-08 13:34:24.387065', 1);
INSERT INTO account VALUES (31, '', 'player222', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-12 00:00:00', '', '2017-04-12 13:29:13.367934', '2017-04-12 13:29:13.367934', '2017-04-12', '', '', '', '', 0, '2017-04-12 13:29:13.367934', 2);
INSERT INTO account VALUES (19, '', 'player888', '01e6a06fe0d37c59d3d75f0806db191f7f2def3fce310ce836977d07604af279', '', '', '', NULL, '', '', 0, '', '', '', '2017-04-08 00:00:00', '::1', '2017-05-06 14:05:36.391733', '2017-04-08 09:47:46.560825', '2017-04-08', '', '', '', '', 0, '2017-04-12 13:10:33.75234', 2);


--
-- Name: account_account_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('account_account_id_seq', 36, true);


--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO cart VALUES ('sfsdfsf                                                                                                                         ', 21, 0, true, '2017-04-28 17:53:51.240452');
INSERT INTO cart VALUES ('sfsdfsf                                                                                                                         ', 22, 0, true, '2017-04-28 17:53:51.240452');
INSERT INTO cart VALUES ('sfsdfsf                                                                                                                         ', 23, 0, true, '2017-04-28 17:53:51.240452');
INSERT INTO cart VALUES ('                                                                                                                                ', 24, 1, true, '2017-04-28 14:37:30.438761');
INSERT INTO cart VALUES ('                                                                                                                                ', 30, 5, true, '2017-04-28 14:36:50.890947');
INSERT INTO cart VALUES ('                                                                                                                                ', 27, 7, true, '2017-04-28 15:08:55.865619');
INSERT INTO cart VALUES ('                                                                                                                                ', 29, 1, true, '2017-04-28 15:38:03.658555');
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
INSERT INTO cart VALUES ('898e181cfb20842a11000bf593ba9a841a5bf714a2f60e339d4ba850c8fb562a                                                                ', 26, 1, true, '2017-04-29 16:14:53.86801');
INSERT INTO cart VALUES ('898e181cfb20842a11000bf593ba9a841a5bf714a2f60e339d4ba850c8fb562a                                                                ', 25, 7, true, '2017-04-29 16:14:53.86801');


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO category VALUES (13, 140, 'Food', '', true, '{/assets/uploads/roxy-fileman-logo.gif}', 3, NULL, '');
INSERT INTO category VALUES (9, 9225, 'Computers', '<p><img alt="" src="http://localhost:3001/assets/uploads/1e43d48s-960.jpg" style="height:533px; width:800px" /></p>
', true, '{/assets/uploads/1e43d48s-960.jpg}', 0, NULL, '');
INSERT INTO category VALUES (10, 18025, 'Laptops', '<p><img alt="" src="http://localhost:3001/assets/uploads/roxy-fileman-logo.gif" style="height:127px; width:289px" /></p>
', true, '{/assets/uploads/pic5.jpg,/assets/uploads/pic7.jpg}', 0, NULL, '');
INSERT INTO category VALUES (23, 500, 'Animals pets', '', true, '{/assets/uploads/Documents/uugugugu.jpg,/assets/uploads/Documents/000006801x.jpg}', 6, NULL, '');
INSERT INTO category VALUES (11, 150, 'Furniture', '', true, '{/assets/uploads/1e43d48s-960.jpg}', 0, NULL, '');
INSERT INTO category VALUES (19, 10025, 'Domestics', '', true, '{}', 0, NULL, '');
INSERT INTO category VALUES (20, 100, 'Notebooks', '<p><em><strong>Washing machines</strong></em></p>

<p><em><strong><img alt="" src="http://localhost:3001/assets/uploads/DSC_2987.jpg" style="height:469px; width:700px" /></strong></em></p>
', true, '{/assets/uploads/1e43d48s-960.jpg,/assets/uploads/21.gif}', 2, 10, '');
INSERT INTO category VALUES (21, 100, 'Ultrabooks', '', true, '{}', 1, 20, '');
INSERT INTO category VALUES (22, 100, 'Tablets', '', true, '{/assets/uploads/bag.jpg,/assets/uploads/bo.jpg}', 0, 21, '');
INSERT INTO category VALUES (24, 101, 'Transport', '', true, '{/assets/uploads/images.jpg,"/assets/uploads/скачанные файлы.jpg",/assets/uploads/21.gif}', 0, NULL, '');


--
-- Name: category_category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('category_category_id_seq', 24, true);


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

INSERT INTO img VALUES (1, 'boss', '2017-04-03 22:21:26.034566');
INSERT INTO img VALUES (2, 'boss', '2017-04-03 22:21:26.034566');


--
-- Name: img_img_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('img_img_id_seq', 2, true);


--
-- Data for Name: login_attempts; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: login_attempts_attempt_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('login_attempts_attempt_id_seq', 29, true);


--
-- Data for Name: permission; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO permission VALUES (8, 'admin', 3);
INSERT INTO permission VALUES (9, 'admin', 4);
INSERT INTO permission VALUES (10, 'admin', 2);


--
-- Name: permission_permission_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('permission_permission_id_seq', 10, true);


--
-- Data for Name: position; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "position" VALUES (1, 'user', 2, 100, true);
INSERT INTO "position" VALUES (2, 'admin', 4, 100, true);
INSERT INTO "position" VALUES (3, 'boss', 0, 100, true);
INSERT INTO "position" VALUES (4, 'manager', 3, 100, true);


--
-- Name: position_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('position_position_id_seq', 5, true);


--
-- Data for Name: product; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO product VALUES (25, 25684.18, '554 Some pet dddd', '', '2017-04-26 07:52:57.60257', '2017-04-29 17:55:59.669302', true, '''554'':1A ''dddd'':4A ''pet'':3A', '{/assets/uploads/pic11.jpg,/assets/uploads/pic2.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (23, 235.00, 'Some pet good', '<h3><var><em><strong>Возвращаясь к пункту про структуру директории, gb трактует всё, что находится в&nbsp;src/, как код вашего </strong></em></var></h3>

<hr />
<h3><var><em><strong>проекта. Все зависимые пакаджи устанавливаются в директор</strong></em></var></h3>

<hr />
<h3><var><em><strong>ию&nbsp;vendor/&nbsp;и именно оттуда код берется при сборке с помощью gb.</strong></em></var></h3>
', '2017-04-26 07:52:43.79536', '2017-05-01 15:05:57.121359', true, '''gb'':10,38 ''good'':3A ''pet'':2A ''src'':16 ''vendor'':28 ''берет'':33 ''ваш'':19 ''возвра'':4 ''всё'':12 ''директор'':9,26 ''зависим'':22 ''и'':27 ''имен'':30 ''код'':18,32 ''наход'':14 ''оттуд'':31 ''пакадж'':23 ''помощ'':37 ''проект'':20 ''пункт'':6 ''сборк'':35 ''структур'':8 ''тракт'':11 ''устанавлива'':24', '{/assets/uploads/ch.jpg,/assets/uploads/pi.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (26, 287.00, '765 Some pet', '', '2017-04-26 07:53:04.817909', '2017-05-01 15:06:24.220413', true, '''765'':1A ''pet'':3A', '{/assets/uploads/pic.jpg,/assets/uploads/pic10.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (21, 102.00, '2558 My Product', '', '2017-04-25 17:48:37.961616', '2017-05-06 13:13:16.332519', true, '''2558'':1A ''product'':3A', '{/assets/uploads/bag.jpg,/assets/uploads/1e43d48s-960.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (29, 647.21, 'Bag ccc', '', '2017-04-27 15:34:40.735855', '2017-05-01 10:43:54.115379', true, '''bag'':1A ''ccc'':2A', '{/assets/uploads/s4.jpg,/assets/uploads/sh.jpg,/assets/uploads/pic4.jpg}', '', 589.36, 0, '{}', 0);
INSERT INTO product VALUES (22, 152.00, 'My ProductOOOKKK', '', '2017-04-26 07:43:57.314561', '2017-05-01 15:05:23.20969', true, '''productoookkk'':2A', '{/assets/uploads/21.gif,/assets/uploads/11395_original.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (24, 186.00, '344 Some pet', '', '2017-04-26 07:52:50.973054', '2017-05-01 15:05:39.843489', true, '''344'':1A ''pet'':3A', '{/assets/uploads/roxy-fileman-logo.gif,/assets/uploads/000006801x.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (30, 312.00, 'Girls ccc', '', '2017-04-27 15:35:17.124572', '2017-05-01 15:06:40.654448', true, '''ccc'':2A ''girl'':1A', '{/assets/uploads/pic12.jpg,/assets/uploads/pi4.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (27, 339.00, 'Shoes nike', '', '2017-04-27 15:33:40.925847', '2017-05-01 15:06:59.854226', true, '''nike'':2A ''shoe'':1A', '{/assets/uploads/sh.jpg,/assets/uploads/s4.jpg}', '', 0.00, 0, '{}', 0);
INSERT INTO product VALUES (28, 378.00, 'Parfume ccc', '', '2017-04-27 15:34:09.240468', '2017-05-01 15:07:17.379328', true, '''ccc'':2A ''parfum'':1A', '{/assets/uploads/pic4.jpg,/assets/uploads/s4.jpg,/assets/uploads/sh.jpg}', '', 0.00, 0, '{}', 0);
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


--
-- Name: product_product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('product_product_id_seq', 31, true);


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO sessions VALUES ('3fb809a2827113254fc9c1212d9d05fb592ce7922b605326a88a8f9547373bf0', '::1', 1492006515, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('40fcf1abb05c2919b0509f96f52c245c281420c612f21b7aef8c3231ca7251df', '::1', 1492097534, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('406107292ccae69675f18d90faa6357a049f61188c43f894bdad6cec16c4e3c7', '::1', 1492191524, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('3a3e71d921e980212844fa64e031db81a9067f87339021aa5735c98aff284e74', '::1', 1492317304, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('67a3cea0b329a294c5a00ca257eb96c29db2b921a86c3d5dc456d66354a955c0', '::1', 1492324783, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('2535d07d02492d6c3dc95f2d191bc29b964be0df0333b12fe6a3f9f26f38785b', '::1', 1492325556, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.138 Safari/537.36 Vivaldi/1.8.770.54', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('630a125d67937f18b224ad0f755afa1a88be13e4f7a47f29e728f34796d82490', '::1', 1492330032, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('5019e03fbfe112b05fffe249288ebc6246cc9c8cdc9bc1ccf84fc9a91d7227e0', '::1', 1492446780, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('365e4b43f6e3bed8b6193cbf10e665cc2985f240e3b7318bd34163ef6f882e27', '::1', 1492517755, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('8444f26bdd0b0e0a6177f6d61798eee06b8b2edd133959e09534ef57865bd078', '::1', 1492522251, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('bb1d67419345c1179f317de91d546902760c159f1ed69a8e046fcc4916198927', '::1', 1492526496, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('a7d3c3102d30faf7693a2489174cdfdbf86dbfb18b2642596c5642e0be844fad', '::1', 1492530509, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('c69ca65deb25a099e8c24637398d3ecd29c090baa8d3018155bad9865847d934', '::1', 1492535065, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('ad8cc8f6ec3b1a84261fd883575d3964c6d6ff015a2e09fe6a5a1ef057e1d4cb', '::1', 1492540658, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('8990b827486b1755691c54ab482506bc4e617704367b5f1844d657e21ed1a225', '::1', 1492610970, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('1ab2346556f5dcaae7e2bb9e50a79fbd25f7ce57cfd26d209be971799034378d', '::1', 1492614590, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('b9c62c9e168def0ee2465bf4a7e4bcfa61e2478ba5c25615da235d09f407decf', '::1', 1492621021, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('56a9e18ec1768125b0a1dff4677d6b9b93fb14fedc9b8d42021684b21c864c21', '::1', 1492672923, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('e33079f744055655ba4eb1c7318b2471ad46198d354852e7afddf79d02571f18', '::1', 1492677773, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('5f1cb535ac6e4e7d0743a033d5e49e8cd8049fc726ab3ff68f8a5d894adaaf0a', '::1', 1492681570, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('152f3f8c38663a69f4d918bb120968c835dd08dbb33c114694415cb5d1787869', '::1', 1492696710, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('93e4533d97c7bb5b49d013a3c989c01cf323d04809fa0ff8b5658c14e2a6921f', '::1', 1492700932, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('ce58c88c6218936d4dbc7d354fbdf22d629443420d82640fe165332bf3cebb1f', '::1', 1492705162, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('784e89bf3394fa9e84fc1b007104248b7cd4471662ee65889409cf3de1ffdd39', '::1', 1492761699, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('decc4553d0740123fb9b58a8fa7fe40da063ded715dbd2e0093606fe3ba85bc5', '::1', 1492770959, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('43ed9f0708a72c3b5cfdbb7e552a22ac17e034d426fad7a1d3024e8f535d16ad', '::1', 1492780013, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('d76f224a0edcee6d73fa4cd460a9bdfcff76690d30f82c2693bcb0e77e6436bb', '::1', 1492783818, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('40722157fae184e64bc3fc19bceb98b33b32d1dc2e6447854906c12a8a7580cd', '::1', 1492794601, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('d856ee1c28b44cf236deb1eaaa62b4427c054a304b595f297347411f08d7d0c3', '::1', 1492844746, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('66524ce9f4cc51a2e44981066eb01a354d9d45211733b79925a04b17fd3b4a82', '::1', 1492850179, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('6c5064c034a413e4a580d09ccf9c1f10de025b3ab8b1bfcb876c27d7bdfea1e8', '::1', 1492853785, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('eb9be8322f5052ff39a57b2500524bfd0179c55b656dcf728b01b44915ca55b5', '::1', 1492858668, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('310f87886e78cb473b08967f76d71a051207a2a2030832b47758c57f646a0655', '::1', 1492873180, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('0a29bd6a2a7db449b77b6fcbf552c2b745935f310cde485fb9af40b43995232f', '::1', 1492876840, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('f5d0ab280cabba3ee97ebdeeb9510c32f3a9b50da10967e7e8288e744cea276c', '::1', 1492880605, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('29464f3038c3d971379554ab297501107a411d631b4037f26d55de77babb1551', '::1', 1492884389, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('ce1997f9c1deabef23fc77f15bd52970bdaf1637feec173a5c4116abdf6528b2', '::1', 1492931549, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('d1be43e842007cff59d2b0dba1a10d4cdc2c17125d704b1461aa18799acb90aa', '::1', 1492950311, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('e6a9f7473f71587e0860e271d4e54599ef45f3deb993c15e7704f367fe675188', '::1', 1492954758, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('d7d94f82288116c786723a0bc482edae5a8add3fdc0649808e7918caf95b935e', '::1', 1492960246, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('551caf49fd29bd7f86b9ebce6e3f224035d920b1b2ae49527ffa61122b79a5fd', '::1', 1492966195, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('0c59735f4c8ed698010963b019aa8a0a770d8344fdfe1452de4e9fee8c029efc', '::1', 1492970169, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('dac055b4a6f9f647ee9babd479d39318839ea1f485a70bafdf750aab2c21e362', '::1', 1493016099, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('b152c53509ae32f79fe83ad5d7b4f0b12a0f8756eda0a490b36cc5cb2c53e1e9', '::1', 1493019856, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('8156bc238758136858a14801101e724c3cf330d988b6ae778586133f5723c31e', '::1', 1493023758, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('eb14b2edc64fa8210f94403d17f91989e8982dbf48abd40bebf26ae8fbc8b65d', '::1', 1493031924, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('ad447e499c2f7ea921203021b3bf9538488086e71dfc9b6a11bfb38180b4e893', '::1', 1493035538, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('c634861ed84fac7d1bbf52926336209f761bdf814a4a2c8a08aa1814bc3ca17d', '::1', 1493039480, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('1e7b2c6fa7c76467ebed0620e5f03325122b8b08f0dc09ce46502ac209978a96', '::1', 1493043919, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('fe07e83229d56008e7ea452955de17562b8e7eb29c9cfbdebc378f15dcddc8ed', '::1', 1493104546, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('7250f8a3fc15c557e1e5f093471bba4b23c29d9d0bc4455d9b3428b3bc0010ce', '::1', 1493108188, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('247cdf2155435327f7022e44cfc8ecaa7d7399aef6337f079fbb9a341332e044', '::1', 1493111801, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('951cf063f2458fc4139717e895eae5a88130c05f751938b9bceb1d920527d46f', '::1', 1493115554, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('c6403e182a47628ffe0531d0f4c9c7248ac68aabbfb3447313a73a5331334852', '::1', 1493140011, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('41456c29dae7a2ee59a3e45f28d2dfabae82812506fecc40d962d06c2c224f10', '::1', 1493187544, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('26aace808f623ded54b133ebca39a668258212665a61084e3b525e804e50a3a1', '::1', 1493191877, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b6f94ea6f3f920b9c1e6c4d285ecf047fddf487293f1e8dbd523e55bb739920b', '::1', 1493191881, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('da76f8c80e1ae3b3a04ec458f3404c93799e32626e5075fc6a718761c748eee9', '::1', 1493208289, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('b32c3fb9db62b25432cae87d092bcd46d7dff883d527d944ffb9bf2c165e0a66', '::1', 1493216310, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('478752710fc6ebbcc221da065708170d4e94567a9c6adf38f2af5a9ddd7ab90c', '::1', 1493220412, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('706647432acde93da9ae5d8bb71c8bf7bfeafda7e8244f967f58bdcd55d0e7ec', '::1', 1493228838, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('4ffb36bca52662e248fa4550f313422ff7ddb1b2b5a88c3219132f1d03c4ce05', '::1', 1493306302, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('b91e08f9d98329c220be65a371362414444032f66d3ad296df345cc54905accb', '::1', 1493310158, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('eb770ce1602a8101881cd5bf5ac11594872d8febdc2d096b06dab80cbd013183', '::1', 1493313844, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b74d9e54740b5f1cf1f9f57085ad2cbf516d2b09ccc1138f196992bf2763d29c', '::1', 1493368749, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b8c4a29784274521fde39e90be277ed749a34acfd0e3ce1d6b0fb1003f1f2059', '::1', 1493388505, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('a97e4efc5d561f2f533701b5ea5e9f746d14d7f7086b5f2fb7d0555e4444de86', '::1', 1493392131, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('029467c99bcf0588ed606a590af594d07d62e7bfab046da46a7c3b63f22f2d58', '::1', 1493396043, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('d64f7fbbb565b0e6da9020aa2193e0c35ca813e542a22d729c3a4db1438dda50', '::1', 1493401295, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('e94995df7849ff4884a155cc73e4971aa7a625b919f958b3e6bf0e2da1f9a713', '::1', 1493450960, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('d671aa47abc9e1b6fcf0a2f693daf6a1abb4fdc0304900770683e0a724702837', '::1', 1493460361, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('a2c40e9ad2b81c76d3247cd8d170ddd663c3b3b65b485b3e60fcbfa4a3607c2f', '::1', 1493465922, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('6bbbaa29cc4579cdef541240b48f1fcc2548fb218e5ddccf97b48236d86bf492', '::1', 1493469709, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('e4b014e8b4c633bebc6f4a4cc7a18f341a0de2f068f18c4ef071ed754c41c59b', '::1', 1493473334, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('46973993e7805271b2ee386e9ceca556fb0d0ecf5e28772f37d99fe95f7ae3a2', '::1', 1493476994, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('156c5530136e9fcdf9c080d3529c805ed177a10c53c86b60fe4a03759731d83b', '::1', 1493482109, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('5e8d7f2299626e443b5c188aedd3885ade3c964134833dece13bd21cf6d17439', '::1', 1493485913, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('da60257f8816c1f6c979620c45c58e5acc73a00be92fac9f918e5618b3f119cc', '::1', 1493548969, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('f347c44b6bbee503e9309f6beb66543f75ac04110eeda9654e31ef6c30488c38', '::1', 1493560960, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('92a655daa8a755ed83e5a1bb3bc117e425b5f4881be0145298c40961075c989d', '::1', 1493626087, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('5c049e70bea629283f3767f05a95d9170f086600aeb649a20f6bfcc68fdcbac8', '::1', 1493630690, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('f4d19b687d9d05f529e6f19c27d5e0886e6c095c417002af19b5dc33f22a349b', '::1', 1493634518, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('07057861e127326a8e25bd0d0b5e480130e939b3088efae6c76dd2566ca35fce', '::1', 1493640289, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('6e55e9610bedde01623a1eccb6e20645b4acd57b959f663f853d3e50df5f94ce', '::1', 1493640290, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('59167e210079b2f983a19c22dd553492951c3cff6ed25ae497789b30e506a50a', '::1', 1493640292, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b6e5e60ebc7cffc83af356e2fd04d27db15e619f442138c3997924c2af75ef1c', '::1', 1493640297, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('18d9cf70d0780b998d19776bec0fc09f8bc6f338578d053c54a72144e62d64e6', '::1', 1493640327, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('ee38446d89d968101e0adc62564aa7088099384e8c62ad8f234db5e5393d5083', '::1', 1493640387, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('068adcfe801154dcd803526a9547033dca9e0eede76523619baaf53bb6a51039', '::1', 1493642005, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('d28311ff6cd272227d1cb1a470d94921a9b8524d9107832b3d81a80131c9207d', '::1', 1493650401, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('2f38f9dd323c32937d065f0f9858026e984e1d49c18ff2bbf9bc97dffa0f3052', '::1', 1493833913, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b5ac13dd07692373152cd9a6690600c6170ac212e3bfddc64b563979808acbe3', '::1', 1493915110, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('b9fadab9a5bd4bb944cadfdf695a93071738c2a5ca91991c43c88d1f71168377', '::1', 1493919368, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('00f061b3a25097ec295165b4822c405c52dd9e88654ee553f294e325048a0501', '::1', 1493976063, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('0e0eb18836658caedefdcb9ea51d924ec0b53adc8d28bc5b537ac51fd4ae7713', '::1', 1494001147, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('a1294b18f0878cca01e4dac540dd272195f462e7aa4d358786e548da92ae32a9', '::1', 1494006096, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('e3fc3013f0c465195bca2877f97b199e6e89b6450c18d21c8613d226a8adc49f', '::1', 1494006097, 0, '', '', 'unauth', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 0, '', '', 0, '');
INSERT INTO sessions VALUES ('86535112b8fe01420abe298d6fffcb9204059448c1dcd1b8392949caee20e808', '::1', 1494006696, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('291dc27725929655b1485f0872dce6dd5d8d8c95b64f171e210cc719313217e2', '::1', 1494055652, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('676b1e69d891e4df7b5f9e88b0054a02a3cea1841a0650dac9546c5caa687ac8', '::1', 1494074116, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');
INSERT INTO sessions VALUES ('3663d84f65bf40b72fd995fb7976750630532a82125a0cbeadfec72f2a6fb718', '::1', 1494078521, 19, '', 'player888', 'logged', 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36', false, 2, 'admin', 'admin', 0, '');


--
-- Data for Name: shipment; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: shipment_shipment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('shipment_shipment_id_seq', 1, false);


--
-- Data for Name: user_temp; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO user_temp VALUES (9, '', 'player5', '123456', '312cde9199032d2138d64a68a85ddd6c0494a5e6a7c338b935bb2491c8eae4c4', '::1', 1491583954, '');


--
-- Name: user_temp_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('user_temp_id_seq', 10, true);


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
-- Name: login_attempts login_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY login_attempts
    ADD CONSTRAINT login_attempts_pkey PRIMARY KEY (attempt_id);


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
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: shipment shipment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY shipment
    ADD CONSTRAINT shipment_pkey PRIMARY KEY (shipment_id);


--
-- Name: user_temp user_temp_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_temp
    ADD CONSTRAINT user_temp_pkey PRIMARY KEY (id);


--
-- Name: product update_product_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_product_updated BEFORE UPDATE ON product FOR EACH ROW EXECUTE PROCEDURE update_product_updated_column();


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
    ADD CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(category_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: permission fk_permission_position_position_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission
    ADD CONSTRAINT fk_permission_position_position_id FOREIGN KEY (permission_position) REFERENCES "position"(position_id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: product_category fk_product_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product_category
    ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES product(product_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

