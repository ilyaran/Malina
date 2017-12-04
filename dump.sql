--
-- PostgreSQL database dump
--

-- Dumped from database version 10.0
-- Dumped by pg_dump version 10.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

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
    product_category,
    product_title,
    product_description,
    product_code,
    product_price,
    product_price1,
    product_enable,
    product_img,
    search_vector)
  VALUES
    (inCategoryId,inTitle, inDescription,inCode,inPrice,inPrice1,inEnable,inImgId,
     (setweight(to_tsvector(inTitle), 'A') || to_tsvector(inDescription)));
  SELECT INTO lastInsertId currval('product_product_id_seq');

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

 /* UPDATE category SET
    (category_quantity)
    =
    coalesce((
      SELECT category_quantity-1 FROM category
      WHERE category_id =
            (
              SELECT product_category FROM product
              WHERE product_id = inId AND product_category != inCategoryId
            )
    ),0);*/

  UPDATE product SET (
    product_category,
    product_title,
    product_description,
    product_code,
    product_price,
    product_price1,
    product_enable,
    product_img,
    product_updated,
    search_vector)
  =
  (inCategoryId,inTitle, inDescription,inCode,inPrice,inPrice1,inEnable,inImgId,now(),
   (setweight(to_tsvector(inTitle), 'A') || to_tsvector(inDescription)))
  WHERE product_id = inId;

  /*UPDATE category SET
    (category_quantity)
    =
    coalesce((
      SELECT category_quantity+1 FROM category
      WHERE category_id = (
              SELECT product_category FROM product
              WHERE product_id = inId AND product_category != inCategoryId
      ) AND category_id = inCategoryId
    ),0);*/

  RETURN 1;
END;
$_$;


ALTER FUNCTION public.product_update(bigint, bigint, character varying, character varying, character varying, numeric, numeric, boolean, character varying[]) OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE account (
    account_id bigint NOT NULL,
    account_email character varying(512),
    account_login character varying(128),
    account_phone character varying(64),
    account_token character varying(512),
    account_last_logged timestamp without time zone DEFAULT now() NOT NULL,
    account_last_ip character varying(128),
    account_updated timestamp without time zone DEFAULT now() NOT NULL,
    account_password character varying(1024),
    account_ban boolean DEFAULT false NOT NULL,
    account_reason character varying(1024) DEFAULT ''::character varying NOT NULL,
    account_ban_duration bigint,
    account_created timestamp without time zone DEFAULT now() NOT NULL,
    account_role bigint
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
    category_parent bigint,
    category_title character varying(255) DEFAULT ''::character varying NOT NULL,
    category_description text,
    category_sort bigint DEFAULT 100 NOT NULL,
    category_created timestamp without time zone DEFAULT now() NOT NULL,
    category_updated timestamp without time zone DEFAULT now() NOT NULL,
    category_quantity bigint DEFAULT 0 NOT NULL,
    category_img character varying(255)[],
    category_parameter bigint[],
    category_enable boolean DEFAULT true NOT NULL,
    category_lang character varying(24) DEFAULT ''::character varying NOT NULL
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
-- Name: news; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE news (
    news_id bigint NOT NULL,
    news_title character varying(255) DEFAULT ''::character varying NOT NULL,
    news_description text,
    news_created timestamp without time zone DEFAULT now() NOT NULL,
    news_updated timestamp without time zone DEFAULT now() NOT NULL,
    news_img character varying(255)[],
    news_account bigint,
    news_comment bigint[],
    news_enable boolean DEFAULT true NOT NULL,
    news_category bigint,
    news_views bigint DEFAULT 0 NOT NULL,
    news_like bigint DEFAULT 0 NOT NULL,
    news_short_description character varying(1024) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE news OWNER TO postgres;

--
-- Name: news_news_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE news_news_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE news_news_id_seq OWNER TO postgres;

--
-- Name: news_news_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE news_news_id_seq OWNED BY news.news_id;


--
-- Name: parameter; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE parameter (
    parameter_id bigint NOT NULL,
    parameter_title character varying(255) DEFAULT ''::character varying NOT NULL,
    parameter_parent bigint,
    parameter_sort bigint DEFAULT 100 NOT NULL,
    parameter_description text,
    parameter_value character varying(1024) DEFAULT ''::character varying NOT NULL,
    parameter_created timestamp without time zone DEFAULT now() NOT NULL,
    parameter_updated timestamp without time zone DEFAULT now() NOT NULL,
    parameter_enable boolean DEFAULT true NOT NULL
);


ALTER TABLE parameter OWNER TO postgres;

--
-- Name: parameter_parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE parameter_parameter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE parameter_parameter_id_seq OWNER TO postgres;

--
-- Name: parameter_parameter_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE parameter_parameter_id_seq OWNED BY parameter.parameter_id;


--
-- Name: permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE permission (
    permission_id bigint NOT NULL,
    permission_title character varying(255) DEFAULT ''::character varying NOT NULL,
    permission_description text,
    permission_enable boolean DEFAULT true NOT NULL
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
-- Name: product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE product (
    product_id bigint NOT NULL,
    product_code character varying(512) DEFAULT ''::character varying NOT NULL,
    product_category bigint,
    product_parameter bigint[],
    product_img character varying(255)[],
    product_title character varying(255) DEFAULT ''::character varying NOT NULL,
    product_description text,
    product_price numeric(20,2) DEFAULT 0.00 NOT NULL,
    product_price1 numeric(20,2) DEFAULT 0.00 NOT NULL,
    product_price2 numeric(20,2) DEFAULT 0.00,
    product_quantity numeric DEFAULT 0 NOT NULL,
    product_sold bigint DEFAULT 0,
    product_views bigint DEFAULT 0 NOT NULL,
    product_comment bigint[],
    product_created timestamp without time zone DEFAULT now() NOT NULL,
    product_updated timestamp without time zone DEFAULT now() NOT NULL,
    product_enable boolean DEFAULT true NOT NULL,
    search_vector tsvector,
    product_like bigint DEFAULT 0 NOT NULL,
    product_star character(1) DEFAULT '1'::bpchar NOT NULL,
    product_short_description character varying(2048) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE product OWNER TO postgres;

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
-- Name: role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE role (
    role_id bigint NOT NULL,
    role_parent bigint,
    role_sort bigint DEFAULT 100 NOT NULL,
    role_title character varying(255) DEFAULT ''::character varying NOT NULL,
    role_enable boolean DEFAULT true NOT NULL,
    role_permission bigint[]
);


ALTER TABLE role OWNER TO postgres;

--
-- Name: role_role_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE role_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE role_role_id_seq OWNER TO postgres;

--
-- Name: role_role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE role_role_id_seq OWNED BY role.role_id;


--
-- Name: session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE session (
    session_id character varying(512) NOT NULL,
    session_account bigint,
    session_created bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    session_data character(1)[],
    session_user_agent character varying(512) DEFAULT ''::character varying NOT NULL,
    session_ip character varying(128) DEFAULT ''::character varying NOT NULL,
    session_email character varying(512),
    session_login character varying(128),
    session_token character varying(512),
    session_phone character varying(64),
    session_device character varying(512) DEFAULT ''::character varying NOT NULL,
    session_role bigint
);


ALTER TABLE session OWNER TO postgres;

--
-- Name: account account_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account ALTER COLUMN account_id SET DEFAULT nextval('account_account_id_seq'::regclass);


--
-- Name: category category_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY category ALTER COLUMN category_id SET DEFAULT nextval('category_category_id_seq'::regclass);


--
-- Name: news news_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY news ALTER COLUMN news_id SET DEFAULT nextval('news_news_id_seq'::regclass);


--
-- Name: parameter parameter_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY parameter ALTER COLUMN parameter_id SET DEFAULT nextval('parameter_parameter_id_seq'::regclass);


--
-- Name: permission permission_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission ALTER COLUMN permission_id SET DEFAULT nextval('permission_permission_id_seq'::regclass);


--
-- Name: product product_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product ALTER COLUMN product_id SET DEFAULT nextval('product_product_id_seq'::regclass);


--
-- Name: role role_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY role ALTER COLUMN role_id SET DEFAULT nextval('role_role_id_seq'::regclass);


--
-- Data for Name: account; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO category VALUES (3, 1, 'Desktop', '', 120, '2017-11-14 21:20:37.199396', '2017-11-14 21:20:37.199396', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (6, 2, 'Ultrabooks', '', 110, '2017-11-15 14:41:56.150225', '2017-11-15 14:41:56.150225', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (7, 2, 'Notebooks', '', 100, '2017-11-15 14:44:06.041425', '2017-11-15 14:44:06.041425', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (5, 4, 'Smartphones', '', 100, '2017-11-15 14:41:26.656251', '2017-11-15 14:41:26.656251', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (4, NULL, 'Gadgets', '', 110, '2017-11-15 14:40:44.758485', '2017-11-15 14:40:44.758485', 2, NULL, NULL, true, '');
INSERT INTO category VALUES (13, NULL, 'Displays', '&lt;p&gt;Letter | 15 November 2017&lt;/p&gt;

&lt;h3&gt;&lt;a href=&#34;https://www.nature.com/articles/nature24649&#34;&gt;&lt;img alt=&#34;&#34; src=&#34;https://media.springernature.com/w75h75/nature-static/assets/v1/image-assets/nature24649-f1.jpg&#34; style=&#34;height:75px; width:75px&#34; /&gt;PD-1 is a haploinsufficient suppressor of T cell lymphomagenesis&lt;/a&gt;&lt;/h3&gt;

&lt;p&gt;Loss of the PD-1 receptor promotes the development of T cell non-Hodgkin lymphomas by modulating oncogenic signalling… show more&lt;/p&gt;

&lt;ul&gt;
	&lt;li&gt;Tim Wartewig&lt;/li&gt;
	&lt;li&gt;, Zsuzsanna Kurgyis&lt;/li&gt;
	&lt;li&gt;&lt;a href=&#34;javascript:;&#34; title=&#34;Show all 12 authors&#34;&gt;[…]&lt;/a&gt;&lt;/li&gt;
	&lt;li&gt;Jürgen Ruland&lt;/li&gt;
&lt;/ul&gt;
', 150, '2017-11-16 13:03:38.596132', '2017-11-16 13:03:38.596132', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (15, NULL, 'Network ', '', 170, '2017-11-16 13:20:56.537135', '2017-11-16 13:20:56.537135', 0, NULL, NULL, true, '');
INSERT INTO category VALUES (2, 1, 'Laptops', '<p>News and Views&nbsp;|&nbsp;15 November 2017</p>

<h3><a href="https://www.nature.com/articles/nature24758"><img alt="" src="https://media.springernature.com/w75h75/nature-static/assets/v1/image-assets/nature24758-f1.jpg" style="height:75px; width:75px" />Archaeology: Inequality has deep roots in Eurasia</a></h3>

<p>A study of 64 archaeological sites across four continents shows that the growth of agricultural and political systems&hellip;&nbsp;show more</p>

<ul>
	<li>Michelle Elliott</li>
</ul>
', 100, '2017-10-31 14:36:34.550983', '2017-10-31 14:36:34.550983', 2, '{Asus-VS248H-P-24-LED-LCD-Monitor-16-9-2-ms-P13729418.jpg,1e43d48s-960.jpg}', NULL, true, '');
INSERT INTO category VALUES (1, NULL, 'Computers', '<p>Comps</p>
', 100, '2017-10-31 14:36:34.550983', '2017-10-31 14:36:34.550983', 25, '{8394337.01.prod.jpg}', NULL, true, '');
INSERT INTO category VALUES (16, NULL, 'WiFi modules cc vv', '', 120, '2017-11-19 18:56:41.084862', '2017-11-19 18:56:41.084862', 0, '{"images (45).jpg","images (3).jpg"}', NULL, true, '');


--
-- Data for Name: news; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: parameter; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO parameter VALUES (3, 'Graphic', NULL, 100, '', '', '2017-11-13 06:44:44.790707', '2017-11-13 06:44:44.790707', true);
INSERT INTO parameter VALUES (5, '8 thread', 4, 100, NULL, '', '2017-11-13 07:49:28.44342', '2017-11-13 07:49:28.44342', true);
INSERT INTO parameter VALUES (6, '16 threads', 4, 110, NULL, '', '2017-11-13 07:51:58.725741', '2017-11-13 07:51:58.725741', true);
INSERT INTO parameter VALUES (7, 'Display', NULL, 120, '', '', '2017-11-13 17:38:47.50543', '2017-11-13 17:38:47.50543', true);
INSERT INTO parameter VALUES (8, 'Discrette', 3, 100, '', '', '2017-11-13 21:07:30.633986', '2017-11-13 21:07:30.633986', true);
INSERT INTO parameter VALUES (11, 'retina', 7, 100, '', '', '2017-11-13 21:57:12.411322', '2017-11-13 21:57:12.411322', true);
INSERT INTO parameter VALUES (12, 'IPS matrix dddd', 7, 110, '', '', '2017-11-14 10:41:52.021043', '2017-11-14 10:41:52.021043', true);
INSERT INTO parameter VALUES (20, 'MainBoards', NULL, 200, '', '', '2017-11-19 15:49:33.535599', '2017-11-19 15:49:33.535599', true);
INSERT INTO parameter VALUES (19, 'MainBoards', NULL, 200, '', '', '2017-11-19 15:49:33.42461', '2017-11-19 15:49:33.42461', true);
INSERT INTO parameter VALUES (21, '111MainBoards', NULL, 210, '', '', '2017-11-19 15:52:59.597356', '2017-11-19 15:52:59.597356', true);
INSERT INTO parameter VALUES (24, '333 sdasd dfsdfsdf', NULL, 230, '', '', '2017-11-19 15:56:18.47885', '2017-11-19 15:56:18.47885', true);
INSERT INTO parameter VALUES (25, 'Integrated', 3, 90, '', '', '2017-11-19 18:38:16.143934', '2017-11-19 18:38:16.143934', true);
INSERT INTO parameter VALUES (1, 'Processors gg', NULL, 50, 'tgdrtgdrth fd', '', '2017-10-31 14:41:27.175606', '2017-10-31 14:41:27.175606', true);
INSERT INTO parameter VALUES (2, '4-Core', 1, 60, '&lt;input type=&#34;checkbox&#34;&gt;', '', '2017-10-31 14:41:27.175606', '2017-10-31 14:41:27.175606', true);
INSERT INTO parameter VALUES (4, 'Threads', 2, 70, 'gdrts drfser', '', '2017-11-13 07:48:50.087554', '2017-11-13 07:48:50.087554', true);
INSERT INTO parameter VALUES (10, '2Gb', 9, 100, '', '', '2017-11-13 21:37:13.886115', '2017-11-13 21:37:13.886115', true);
INSERT INTO parameter VALUES (13, '4Gb', 9, 110, '', '', '2017-11-14 11:05:44.317195', '2017-11-14 11:05:44.317195', true);
INSERT INTO parameter VALUES (9, 'GRAM', 8, 100, '', '', '2017-11-13 21:34:45.601092', '2017-11-13 21:34:45.601092', true);
INSERT INTO parameter VALUES (16, '6 cores', 1, 110, '<p>Letter&nbsp;|&nbsp;08 November 2017</p>

<h3><a href="https://www.nature.com/articles/nature24476"><img alt="" src="https://media.springernature.com/w75h75/nature-static/assets/v1/image-assets/nature24476-f1.jpg" style="height:75px; width:75px" />Parallel palaeogenomic transects reveal complex genetic history of early European farmers</a></h3>

<p>In European Neolithic populations, the arrival of farmers prompted admixture with local hunter-gatherers over many&hellip;&nbsp;show more</p>

<ul>
	<li>Mark Lipson</li>
	<li>,&nbsp;Anna Sz&eacute;cs&eacute;nyi-Nagy</li>
	<li><a href="javascript:;" title="Show all 57 authors">[&hellip;]</a></li>
	<li>David Reich</li>
</ul>
', '', '2017-11-16 13:39:28.200614', '2017-11-16 13:39:28.200614', true);


--
-- Data for Name: permission; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO permission VALUES (3, 'product', 'add edit delete read', true);


--
-- Data for Name: product; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO product VALUES (35, 'sdfsdfsdf', 1, NULL, NULL, 'Awesome product sdfsa ASdas', '<p>News&nbsp;|&nbsp;17 November 2017</p>

<h3><a href="https://www.nature.com/news/giant-telescope-s-mobile-phone-dead-zones-rile-south-african-residents-1.22998"><img alt="" src="https://www.nature.com/homepage-assets/npg/news/2017/171114/images/w75h75/nature.2017.22998-i7.47693.jpg" style="height:75px; width:75px" />Giant telescope&rsquo;s mobile-phone &lsquo;dead zones&rsquo; rile South African residents</a></h3>

<p>Sensitive radio dishes of the Square Kilometre Array will affect phone reception &mdash; and could harm local economies, say farmers.</p>
', 250.98, 253.36, 0.00, 0, 0, 0, NULL, '2017-11-19 14:49:35.376475', '2017-11-19 15:32:04.619066', true, '''17'':6 ''2017'':8 ''affect'':30 ''african'':19 ''array'':28 ''asda'':4A ''awesom'':1A ''could'':34 ''dead'':15 ''dish'':23 ''economi'':37 ''farmer'':39 ''giant'':9 ''harm'':35 ''kilometr'':27 ''local'':36 ''mobil'':13 ''mobile-phon'':12 ''news'':5 ''novemb'':7 ''phone'':14,31 ''product'':2A ''radio'':22 ''recept'':32 ''resid'':20 ''rile'':17 ''say'':38 ''sdfsa'':3A ''sensit'':21 ''south'':18 ''squar'':26 ''telescop'':10 ''zone'':16', 0, '1', '');
INSERT INTO product VALUES (34, 'zdfdfgfdgdfg', 1, NULL, '{Asus-VS248H-P-24-LED-LCD-Monitor-16-9-2-ms-P13729418.jpg}', 'Monitor cc', '<p>News and Views&nbsp;|&nbsp;15 November 2017</p>

<h3><a href="https://www.nature.com/articles/nature24760"><img alt="" src="https://media.springernature.com/w75h75/nature-static/assets/v1/image-assets/nature24760-f1.jpg" style="height:75px; width:75px" />Microbiota: A high-pressure situation for bacteria</a></h3>

<p>Analyses in mice suggest that dietary salt increases blood pressure partly by affecting some of the microbes that&hellip;&nbsp;show more</p>

<ul>
	<li>David A. Relman</li>
</ul>
', 345.00, 343.00, 0.00, 0, 0, 0, NULL, '2017-11-16 12:53:03.521017', '2017-11-16 12:53:17.776602', true, '''15'':6 ''2017'':8 ''affect'':29 ''analys'':17 ''bacteria'':16 ''blood'':25 ''cc'':2A ''david'':37 ''dietari'':22 ''high'':12 ''high-pressur'':11 ''increas'':24 ''mice'':19 ''microb'':33 ''microbiota'':9 ''monitor'':1A ''news'':3 ''novemb'':7 ''part'':27 ''pressur'':13,26 ''relman'':39 ''salt'':23 ''show'':35 ''situat'':14 ''suggest'':20 ''view'':5', 0, '1', '');
INSERT INTO product VALUES (32, 'sdfzsdf', 4, NULL, '{"images (10).jpg"}', '112 sdfszd zsdsd', '', 256.36, 254.28, 0.00, 0, 0, 0, NULL, '2017-11-15 20:40:02.331731', '2017-11-15 21:14:07.778512', true, '''112'':1A ''sdfszd'':2A ''zsdsd'':3A', 0, '1', '');


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO session VALUES ('094ad34fdb497994d0c3b535d1c18627e77a2e98da69edb6a85874b0751ccef2', NULL, 1511110208, NULL, 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36', '', NULL, NULL, NULL, NULL, '', NULL);


--
-- Name: account_account_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('account_account_id_seq', 1, false);


--
-- Name: category_category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('category_category_id_seq', 16, true);


--
-- Name: news_news_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('news_news_id_seq', 1, false);


--
-- Name: parameter_parameter_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('parameter_parameter_id_seq', 26, true);


--
-- Name: permission_permission_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('permission_permission_id_seq', 3, true);


--
-- Name: product_product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('product_product_id_seq', 35, true);


--
-- Name: role_role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('role_role_id_seq', 1, false);


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
-- Name: news news_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY news
    ADD CONSTRAINT news_pkey PRIMARY KEY (news_id);


--
-- Name: parameter parameter_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY parameter
    ADD CONSTRAINT parameter_pkey PRIMARY KEY (parameter_id);


--
-- Name: permission permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY permission
    ADD CONSTRAINT permission_pkey PRIMARY KEY (permission_id);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY product
    ADD CONSTRAINT product_pkey PRIMARY KEY (product_id);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY role
    ADD CONSTRAINT role_pkey PRIMARY KEY (role_id);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_pkey PRIMARY KEY (session_id);


--
-- Name: account_account_email_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_email_uindex ON account USING btree (account_email);


--
-- Name: account_account_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_id_uindex ON account USING btree (account_id);


--
-- Name: account_account_login_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_login_uindex ON account USING btree (account_login);


--
-- Name: account_account_phone_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_phone_uindex ON account USING btree (account_phone);


--
-- Name: account_account_token_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX account_account_token_uindex ON account USING btree (account_token);


--
-- Name: category_category_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX category_category_id_uindex ON category USING btree (category_id);


--
-- Name: news_news_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX news_news_id_uindex ON news USING btree (news_id);


--
-- Name: parameter_parameter_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX parameter_parameter_id_uindex ON parameter USING btree (parameter_id);


--
-- Name: permission_permission_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX permission_permission_id_uindex ON permission USING btree (permission_id);


--
-- Name: product_product_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX product_product_id_uindex ON product USING btree (product_id);


--
-- Name: role_role_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX role_role_id_uindex ON role USING btree (role_id);


--
-- Name: session_session_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX session_session_id_uindex ON session USING btree (session_id);


--
-- Name: account account_role_role_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_role_role_id_fk FOREIGN KEY (account_role) REFERENCES role(role_id) ON DELETE SET NULL;


--
-- Name: category category_category_category_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY category
    ADD CONSTRAINT category_category_category_id_fk FOREIGN KEY (category_parent) REFERENCES category(category_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: news news_account_account_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY news
    ADD CONSTRAINT news_account_account_id_fk FOREIGN KEY (news_account) REFERENCES account(account_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_account_account_email_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_email_fk FOREIGN KEY (session_email) REFERENCES account(account_email) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_account_account_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_id_fk FOREIGN KEY (session_account) REFERENCES account(account_id) ON DELETE CASCADE;


--
-- Name: session session_account_account_login_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_login_fk FOREIGN KEY (session_login) REFERENCES account(account_login) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_account_account_phone_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_phone_fk FOREIGN KEY (session_phone) REFERENCES account(account_phone) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_account_account_token_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_account_account_token_fk FOREIGN KEY (session_token) REFERENCES account(account_token) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: session session_role_role_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY session
    ADD CONSTRAINT session_role_role_id_fk FOREIGN KEY (session_role) REFERENCES role(role_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

