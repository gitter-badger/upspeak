--
-- Name: organizations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE organizations (
    id bigint NOT NULL,
    org_slug character varying NOT NULL,
    org_name text,
    org_primary_contact bigint NOT NULL
);


--
-- Name: team_memberships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE team_memberships (
    team_id bigint NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: teams; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE teams (
    id bigint NOT NULL,
    team_name text NOT NULL,
    team_slug character varying(32) NOT NULL,
    parent_org bigint NOT NULL,
    parent_team bigint
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE users (
    id bigint NOT NULL,
    username character varying(32)[] NOT NULL,
    email_primary text NOT NULL,
    name text,
    created_at timestamp with time zone NOT NULL,
    last_updated_at timestamp with time zone NOT NULL
);


--
-- Name: users email; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT email UNIQUE (email_primary);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: organizations slug; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY organizations
    ADD CONSTRAINT slug UNIQUE (org_slug);


--
-- Name: team_memberships team_memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY team_memberships
    ADD CONSTRAINT team_memberships_pkey PRIMARY KEY (team_id, user_id);


--
-- Name: teams teams_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (id);


--
-- Name: teams unique_team_in_org; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY teams
    ADD CONSTRAINT unique_team_in_org UNIQUE (team_slug, parent_org);


--
-- Name: users username; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT username UNIQUE (username);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: org_slug_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX org_slug_idx ON organizations USING hash (org_slug);


--
-- Name: team_user_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX team_user_idx ON team_memberships USING btree (user_id);


--
-- Name: username_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX username_idx ON users USING hash (username);


--
-- Name: organizations primary_contact; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY organizations
    ADD CONSTRAINT primary_contact FOREIGN KEY (org_primary_contact) REFERENCES users(id);


--
-- Name: team_memberships team; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY team_memberships
    ADD CONSTRAINT team FOREIGN KEY (team_id) REFERENCES teams(id);


--
-- Name: team_memberships user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY team_memberships
    ADD CONSTRAINT "user" FOREIGN KEY (user_id) REFERENCES users(id);
