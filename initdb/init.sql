--
-- PostgreSQL database dump
--

-- Dumped from database version 12.19
-- Dumped by pg_dump version 12.19

-- Started on 2024-06-28 09:53:56 UTC

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- TOC entry 203 (class 1259 OID 24637)
-- Name: sessions; Type: TABLE; Schema: public; Owner: server
--

CREATE TABLE public.sessions (
    session_id character varying(255) NOT NULL,
    player1_id character varying(255) NOT NULL,
    player2_id character varying(255) NOT NULL,
    moves jsonb DEFAULT '[]'::jsonb NOT NULL
);


ALTER TABLE public.sessions OWNER TO server;

--
-- TOC entry 202 (class 1259 OID 24627)
-- Name: users; Type: TABLE; Schema: public; Owner: server
--

CREATE TABLE public.users (
    player_id character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL
);


ALTER TABLE public.users OWNER TO server;

--
-- TOC entry 3035 (class 0 OID 24637)
-- Dependencies: 203
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: server
--

COPY public.sessions (session_id, player1_id, player2_id, moves) FROM stdin;
\.


--
-- TOC entry 3034 (class 0 OID 24627)
-- Dependencies: 202
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: server
--

COPY public.users (player_id, username, password) FROM stdin;
\.


--
-- TOC entry 2905 (class 2606 OID 24645)
-- Name: sessions session_pkey; Type: CONSTRAINT; Schema: public; Owner: server
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT session_pkey PRIMARY KEY (session_id);


--
-- TOC entry 2899 (class 2606 OID 24660)
-- Name: users unique_username; Type: CONSTRAINT; Schema: public; Owner: server
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_username UNIQUE (username);


--
-- TOC entry 2901 (class 2606 OID 24634)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: server
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (player_id);


--
-- TOC entry 2902 (class 1259 OID 24656)
-- Name: idx_player1_id; Type: INDEX; Schema: public; Owner: server
--

CREATE INDEX idx_player1_id ON public.sessions USING btree (player1_id);


--
-- TOC entry 2903 (class 1259 OID 24657)
-- Name: idx_player2_id; Type: INDEX; Schema: public; Owner: server
--

CREATE INDEX idx_player2_id ON public.sessions USING btree (player2_id);


--
-- TOC entry 2906 (class 2606 OID 24646)
-- Name: sessions session_player1_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: server
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT session_player1_id_fkey FOREIGN KEY (player1_id) REFERENCES public.users(player_id);


--
-- TOC entry 2907 (class 2606 OID 24651)
-- Name: sessions session_player2_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: server
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT session_player2_id_fkey FOREIGN KEY (player2_id) REFERENCES public.users(player_id);


-- Completed on 2024-06-28 09:53:56 UTC

--
-- PostgreSQL database dump complete
--
