PGDMP         	                x           CoEpi    12.2    12.2 +    �           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            �           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            �           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            �           1262    16386    CoEpi    DATABASE     e   CREATE DATABASE "CoEpi" WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C' LC_CTYPE = 'C';
    DROP DATABASE "CoEpi";
                postgres    false            �           0    0    DATABASE "CoEpi"    COMMENT     M   COMMENT ON DATABASE "CoEpi" IS 'Server side database for the CoEpi project';
                   postgres    false    3257                        3079    16531 	   uuid-ossp 	   EXTENSION     ?   CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
    DROP EXTENSION "uuid-ossp";
                   false            �           0    0    EXTENSION "uuid-ossp"    COMMENT     W   COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
                        false    2            �            1259    16396    client_bluetooth_devices    TABLE     �   CREATE TABLE public.client_bluetooth_devices (
    client_id bigint NOT NULL,
    bluetooth_id bigint NOT NULL,
    bluetooth_hash_prefix text[] NOT NULL
);
 ,   DROP TABLE public.client_bluetooth_devices;
       public         heap    postgres    false            �           0    0    TABLE client_bluetooth_devices    COMMENT     �   COMMENT ON TABLE public.client_bluetooth_devices IS 'A list of all bluetooth devices that a client recognises as ''self'' - eg. your own phone, laptop(s) watch(es) headphone(s) mice etc.';
          public          postgres    false    205            �           0    0 5   COLUMN client_bluetooth_devices.bluetooth_hash_prefix    COMMENT     �   COMMENT ON COLUMN public.client_bluetooth_devices.bluetooth_hash_prefix IS 'The first 4 digits of the one-way hash of a Bluetooth UUID';
          public          postgres    false    205            �            1259    16394 &   ClientBluetoothDevices_BluetoothID_seq    SEQUENCE     �   ALTER TABLE public.client_bluetooth_devices ALTER COLUMN bluetooth_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."ClientBluetoothDevices_BluetoothID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            public          postgres    false    205            �            1259    16387    clients    TABLE     �   CREATE TABLE public.clients (
    client_id bigint NOT NULL,
    client_uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    begin_dt timestamp with time zone NOT NULL,
    last_update_dt timestamp with time zone NOT NULL
);
    DROP TABLE public.clients;
       public         heap    postgres    false    2            �            1259    16409    Clients_ClientID_seq    SEQUENCE     �   ALTER TABLE public.clients ALTER COLUMN client_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Clients_ClientID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            public          postgres    false    203            �            1259    16506    case_bluetooth_devices    TABLE     �   CREATE TABLE public.case_bluetooth_devices (
    case_id bigint NOT NULL,
    case_bluetooth_id bigint NOT NULL,
    bluetooth_hash text NOT NULL
);
 *   DROP TABLE public.case_bluetooth_devices;
       public         heap    postgres    false            �           0    0    TABLE case_bluetooth_devices    COMMENT     l   COMMENT ON TABLE public.case_bluetooth_devices IS 'A list of all bluetooth devices associated with a case';
          public          postgres    false    210            �            1259    16490    cases    TABLE     �   CREATE TABLE public.cases (
    client_id bigint NOT NULL,
    case_id bigint NOT NULL,
    symptom_onset_time time with time zone NOT NULL,
    likely_exposure_time time with time zone
);
    DROP TABLE public.cases;
       public         heap    postgres    false            �           0    0    TABLE cases    COMMENT     �   COMMENT ON TABLE public.cases IS 'A case record representing a person''s illness.  One person may have multiple case records, for separate illnesses.';
          public          postgres    false    209            �            1259    16477 	   exposures    TABLE        CREATE TABLE public.exposures (
    donor_client_id bigint NOT NULL,
    exposure_id bigint NOT NULL,
    case_id bigint NOT NULL,
    geohash text NOT NULL,
    start_timeband time with time zone NOT NULL,
    end_timeband time with time zone NOT NULL
);
    DROP TABLE public.exposures;
       public         heap    postgres    false            �           0    0    TABLE exposures    COMMENT     \   COMMENT ON TABLE public.exposures IS 'A record of potential infectious exposure of others';
          public          postgres    false    208            �            1259    16411    subscriptions    TABLE     �   CREATE TABLE public.subscriptions (
    client_id bigint NOT NULL,
    subscription_id bigint NOT NULL,
    geohash text[] NOT NULL,
    start_timeband timestamp with time zone NOT NULL,
    end_timeband timestamp with time zone NOT NULL
);
 !   DROP TABLE public.subscriptions;
       public         heap    postgres    false            �           0    0    TABLE subscriptions    COMMENT     �   COMMENT ON TABLE public.subscriptions IS 'A record of all clients requesting to be sent alerts for a given physical area (geohash) and time period.';
          public          postgres    false    207            �          0    16506    case_bluetooth_devices 
   TABLE DATA           \   COPY public.case_bluetooth_devices (case_id, case_bluetooth_id, bluetooth_hash) FROM stdin;
    public          postgres    false    210   �4       �          0    16490    cases 
   TABLE DATA           ]   COPY public.cases (client_id, case_id, symptom_onset_time, likely_exposure_time) FROM stdin;
    public          postgres    false    209   �4       �          0    16396    client_bluetooth_devices 
   TABLE DATA           b   COPY public.client_bluetooth_devices (client_id, bluetooth_id, bluetooth_hash_prefix) FROM stdin;
    public          postgres    false    205   5       �          0    16387    clients 
   TABLE DATA           S   COPY public.clients (client_id, client_uuid, begin_dt, last_update_dt) FROM stdin;
    public          postgres    false    203   $5       �          0    16477 	   exposures 
   TABLE DATA           q   COPY public.exposures (donor_client_id, exposure_id, case_id, geohash, start_timeband, end_timeband) FROM stdin;
    public          postgres    false    208   A5       �          0    16411    subscriptions 
   TABLE DATA           j   COPY public.subscriptions (client_id, subscription_id, geohash, start_timeband, end_timeband) FROM stdin;
    public          postgres    false    207   ^5       �           0    0 &   ClientBluetoothDevices_BluetoothID_seq    SEQUENCE SET     W   SELECT pg_catalog.setval('public."ClientBluetoothDevices_BluetoothID_seq"', 1, false);
          public          postgres    false    204            �           0    0    Clients_ClientID_seq    SEQUENCE SET     E   SELECT pg_catalog.setval('public."Clients_ClientID_seq"', 1, false);
          public          postgres    false    206                       2606    16403 4   client_bluetooth_devices ClientBluetoothDevices_pkey 
   CONSTRAINT     ~   ALTER TABLE ONLY public.client_bluetooth_devices
    ADD CONSTRAINT "ClientBluetoothDevices_pkey" PRIMARY KEY (bluetooth_id);
 `   ALTER TABLE ONLY public.client_bluetooth_devices DROP CONSTRAINT "ClientBluetoothDevices_pkey";
       public            postgres    false    205                       2606    16393 "   clients No duplicated Client UUIDs 
   CONSTRAINT     f   ALTER TABLE ONLY public.clients
    ADD CONSTRAINT "No duplicated Client UUIDs" UNIQUE (client_uuid);
 N   ALTER TABLE ONLY public.clients DROP CONSTRAINT "No duplicated Client UUIDs";
       public            postgres    false    203            !           2606    16418     subscriptions Subscription Index 
   CONSTRAINT     m   ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT "Subscription Index" PRIMARY KEY (subscription_id);
 L   ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT "Subscription Index";
       public            postgres    false    207                       2606    16476 $   subscriptions Verify valid time band    CHECK CONSTRAINT     �   ALTER TABLE public.subscriptions
    ADD CONSTRAINT "Verify valid time band" CHECK ((start_timeband < end_timeband)) NOT VALID;
 K   ALTER TABLE public.subscriptions DROP CONSTRAINT "Verify valid time band";
       public          postgres    false    207    207    207    207            '           2606    16513 2   case_bluetooth_devices case_bluetooth_devices_pkey 
   CONSTRAINT        ALTER TABLE ONLY public.case_bluetooth_devices
    ADD CONSTRAINT case_bluetooth_devices_pkey PRIMARY KEY (case_bluetooth_id);
 \   ALTER TABLE ONLY public.case_bluetooth_devices DROP CONSTRAINT case_bluetooth_devices_pkey;
       public            postgres    false    210            %           2606    16494    cases cases_pkey 
   CONSTRAINT     S   ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_pkey PRIMARY KEY (case_id);
 :   ALTER TABLE ONLY public.cases DROP CONSTRAINT cases_pkey;
       public            postgres    false    209                       2606    16391    clients clients_pkey 
   CONSTRAINT     Y   ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (client_id);
 >   ALTER TABLE ONLY public.clients DROP CONSTRAINT clients_pkey;
       public            postgres    false    203            #           2606    16484    exposures exposures_pkey 
   CONSTRAINT     _   ALTER TABLE ONLY public.exposures
    ADD CONSTRAINT exposures_pkey PRIMARY KEY (exposure_id);
 B   ALTER TABLE ONLY public.exposures DROP CONSTRAINT exposures_pkey;
       public            postgres    false    208            (           2606    16404 !   client_bluetooth_devices ClientID    FK CONSTRAINT     �   ALTER TABLE ONLY public.client_bluetooth_devices
    ADD CONSTRAINT "ClientID" FOREIGN KEY (client_id) REFERENCES public.clients(client_id) ON DELETE CASCADE;
 M   ALTER TABLE ONLY public.client_bluetooth_devices DROP CONSTRAINT "ClientID";
       public          postgres    false    205    3101    203            )           2606    16419 )   subscriptions Subscriptions_ClientID_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT "Subscriptions_ClientID_fkey" FOREIGN KEY (client_id) REFERENCES public.clients(client_id) ON DELETE CASCADE;
 U   ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT "Subscriptions_ClientID_fkey";
       public          postgres    false    3101    207    203            -           2606    16514 "   case_bluetooth_devices case verify    FK CONSTRAINT     �   ALTER TABLE ONLY public.case_bluetooth_devices
    ADD CONSTRAINT "case verify" FOREIGN KEY (case_id) REFERENCES public.cases(case_id) ON DELETE CASCADE;
 N   ALTER TABLE ONLY public.case_bluetooth_devices DROP CONSTRAINT "case verify";
       public          postgres    false    209    210    3109            ,           2606    16519    cases cases_client_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.cases
    ADD CONSTRAINT cases_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.clients(client_id) ON DELETE CASCADE NOT VALID;
 D   ALTER TABLE ONLY public.cases DROP CONSTRAINT cases_client_id_fkey;
       public          postgres    false    203    209    3101            *           2606    16485    exposures client reference    FK CONSTRAINT     �   ALTER TABLE ONLY public.exposures
    ADD CONSTRAINT "client reference" FOREIGN KEY (donor_client_id) REFERENCES public.clients(client_id);
 F   ALTER TABLE ONLY public.exposures DROP CONSTRAINT "client reference";
       public          postgres    false    3101    208    203            +           2606    16524     exposures exposures_case_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.exposures
    ADD CONSTRAINT exposures_case_id_fkey FOREIGN KEY (case_id) REFERENCES public.cases(case_id) ON DELETE CASCADE NOT VALID;
 J   ALTER TABLE ONLY public.exposures DROP CONSTRAINT exposures_case_id_fkey;
       public          postgres    false    3109    208    209            �      x������ � �      �      x������ � �      �      x������ � �      �      x������ � �      �      x������ � �      �      x������ � �     