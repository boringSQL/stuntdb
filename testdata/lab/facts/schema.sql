--
-- PostgreSQL database dump
--

\restrict stuntdb

-- Dumped from database version 18.7 (Debian 18.7-1.pgdg13+5)
-- Dumped by pg_dump version 18.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: coupons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coupons (
    code text NOT NULL,
    discount integer NOT NULL
);
ALTER TABLE ONLY public.coupons ALTER COLUMN code SET (pg_attver=29171888);
ALTER TABLE ONLY public.coupons ALTER COLUMN discount SET (pg_attver=29171888);


--
-- Name: customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    tenant_id integer NOT NULL,
    email text NOT NULL,
    created_at timestamp with time zone NOT NULL
);
ALTER TABLE ONLY public.customers ALTER COLUMN id SET (pg_attver=29145080);
ALTER TABLE ONLY public.customers ALTER COLUMN tenant_id SET (pg_attver=29145080);
ALTER TABLE ONLY public.customers ALTER COLUMN email SET (pg_attver=29145080);
ALTER TABLE ONLY public.customers ALTER COLUMN created_at SET (pg_attver=29145080);


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.order_items (
    order_id bigint NOT NULL,
    sku text NOT NULL,
    qty integer NOT NULL,
    price numeric(10,2) NOT NULL
);
ALTER TABLE ONLY public.order_items ALTER COLUMN order_id SET (pg_attver=29219128);
ALTER TABLE ONLY public.order_items ALTER COLUMN sku SET (pg_attver=29219128);
ALTER TABLE ONLY public.order_items ALTER COLUMN qty SET (pg_attver=29219128);
ALTER TABLE ONLY public.order_items ALTER COLUMN price SET (pg_attver=29219128);


--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    tenant_id integer NOT NULL,
    customer_id bigint NOT NULL,
    status text NOT NULL,
    total numeric(12,2) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    coupon_code text
);
ALTER TABLE ONLY public.orders ALTER COLUMN id SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN tenant_id SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN customer_id SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN status SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN total SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN created_at SET (pg_attver=29187608);
ALTER TABLE ONLY public.orders ALTER COLUMN coupon_code SET (pg_attver=29187608);


--
-- Name: tenants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tenants (
    id integer NOT NULL,
    name text NOT NULL,
    plan text NOT NULL
);
ALTER TABLE ONLY public.tenants ALTER COLUMN id SET (pg_attver=29128144);
ALTER TABLE ONLY public.tenants ALTER COLUMN name SET (pg_attver=29128144);
ALTER TABLE ONLY public.tenants ALTER COLUMN plan SET (pg_attver=29128144);


--
-- Name: coupons coupons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupons
    ADD CONSTRAINT coupons_pkey PRIMARY KEY (code);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: tenants tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT tenants_pkey PRIMARY KEY (id);


--
-- Name: customers_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX customers_email_idx ON public.customers USING btree (email);


--
-- Name: customers_tenant_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX customers_tenant_idx ON public.customers USING btree (tenant_id);


--
-- Name: order_items_order_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX order_items_order_idx ON public.order_items USING btree (order_id);


--
-- Name: order_items_sku_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX order_items_sku_idx ON public.order_items USING btree (sku);


--
-- Name: orders_customer_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX orders_customer_idx ON public.orders USING btree (customer_id);


--
-- Name: orders_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX orders_status_idx ON public.orders USING btree (status);


--
-- Name: orders_tenant_created_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX orders_tenant_created_idx ON public.orders USING btree (tenant_id, created_at);


--
-- Name: customers customers_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: orders orders_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- PostgreSQL database dump complete
--

\unrestrict stuntdb

