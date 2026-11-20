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

--
-- Statistics for Name: coupons; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'coupons',
	'relpages', '109'::integer,
	'reltuples', '20000'::real,
	'relallvisible', '109'::integer,
	'relallfrozen', '109'::integer
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'coupons',
	'attname', 'code',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '6'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{C1,C10177,C10357,C10537,C10717,C10898,C11077,C11257,C11437,C11617,C11798,C11978,C12157,C12337,C12517,C12698,C12878,C13057,C13237,C13417,C13598,C13778,C13958,C14137,C14317,C14498,C14678,C14858,C15037,C15217,C15398,C15578,C15758,C15938,C16117,C16298,C16478,C16658,C16838,C17017,C17198,C17378,C17558,C17738,C17918,C18098,C18278,C18458,C18638,C18818,C18999,C19178,C19358,C19538,C19718,C19899,C2077,C2257,C2437,C2617,C2798,C2978,C3157,C3337,C3517,C3698,C3878,C4057,C4237,C4417,C4598,C4778,C4958,C5137,C5317,C5498,C5678,C5858,C6037,C6217,C6398,C6578,C6758,C6938,C7117,C7298,C7478,C7658,C7838,C8017,C8198,C8378,C8558,C8738,C8918,C9098,C9278,C9458,C9638,C9818,C9999}'::text,
	'correlation', '-0.39336452'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'coupons',
	'attname', 'discount',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '4'::integer,
	'n_distinct', '20'::real,
	'most_common_vals', '{5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24}'::text,
	'most_common_freqs', '{0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05,0.05}'::real[],
	'correlation', '0.05071245'::real
);


--
-- Statistics for Name: customers; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'relpages', '2803'::integer,
	'reltuples', '300000'::real,
	'relallvisible', '2803'::integer,
	'relallfrozen', '2573'::integer
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'attname', 'created_at',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{"2024-01-01 00:19:55.550131+00","2024-01-08 03:47:36.61083+00","2024-01-15 00:47:49.427091+00","2024-01-22 15:06:10.255979+00","2024-01-30 04:01:13.630292+00","2024-02-05 20:32:49.641187+00","2024-02-13 02:38:18.131506+00","2024-02-20 00:24:17.472824+00","2024-02-26 20:19:51.17201+00","2024-03-05 20:37:16.540955+00","2024-03-12 20:53:17.471931+00","2024-03-19 19:33:24.113427+00","2024-03-26 20:48:04.224469+00","2024-04-02 09:13:19.253921+00","2024-04-08 22:36:11.970037+00","2024-04-15 07:46:13.284307+00","2024-04-22 13:38:18.889882+00","2024-04-29 04:28:27.019703+00","2024-05-06 12:20:14.900874+00","2024-05-13 22:17:23.928541+00","2024-05-21 06:28:50.219854+00","2024-05-28 12:06:13.910615+00","2024-06-04 03:02:54.799377+00","2024-06-11 10:05:39.474092+00","2024-06-18 03:57:05.526818+00","2024-06-25 02:00:18.808813+00","2024-07-01 19:39:03.198131+00","2024-07-08 18:54:31.670558+00","2024-07-15 13:56:52.584661+00","2024-07-22 16:47:26.476526+00","2024-07-29 06:26:00.722223+00","2024-08-05 12:14:44.605781+00","2024-08-12 08:21:06.231288+00","2024-08-19 16:40:30.166405+00","2024-08-27 09:19:57.703585+00","2024-09-03 10:02:16.989401+00","2024-09-09 09:56:10.458969+00","2024-09-16 10:46:12.565958+00","2024-09-23 00:17:08.275503+00","2024-09-29 22:06:14.158603+00","2024-10-07 05:43:21.988197+00","2024-10-13 19:41:27.23098+00","2024-10-20 16:41:49.412493+00","2024-10-27 18:33:22.086303+00","2024-11-03 02:54:11.905335+00","2024-11-09 21:17:55.061219+00","2024-11-17 00:29:09.296097+00","2024-11-23 20:11:58.016902+00","2024-11-30 23:30:23.485291+00","2024-12-08 07:46:23.531182+00","2024-12-15 13:52:54.254833+00","2024-12-21 21:45:04.340925+00","2024-12-28 11:20:53.582883+00","2025-01-04 21:00:01.238171+00","2025-01-12 14:16:09.387277+00","2025-01-19 23:59:25.010793+00","2025-01-26 23:13:59.47567+00","2025-02-02 12:35:09.821931+00","2025-02-09 16:28:06.236702+00","2025-02-16 16:40:36.050563+00","2025-02-24 05:41:41.130056+00","2025-03-02 23:45:58.674982+00","2025-03-10 09:57:39.797347+00","2025-03-17 14:20:22.064987+00","2025-03-24 23:01:30.53788+00","2025-03-31 10:08:46.715558+00","2025-04-07 22:01:09.562378+00","2025-04-14 18:17:17.107629+00","2025-04-21 07:24:55.122142+00","2025-04-28 07:57:19.578301+00","2025-05-05 08:29:21.1079+00","2025-05-12 13:12:03.946104+00","2025-05-20 07:55:35.298754+00","2025-05-27 21:32:16.555067+00","2025-06-04 04:56:05.578609+00","2025-06-11 09:17:37.69412+00","2025-06-18 18:24:08.750907+00","2025-06-25 20:36:48.531061+00","2025-07-02 13:33:39.228891+00","2025-07-09 22:12:05.464842+00","2025-07-16 21:17:35.113081+00","2025-07-23 11:07:57.771443+00","2025-07-30 05:12:22.100887+00","2025-08-06 04:29:39.776696+00","2025-08-12 10:23:33.795324+00","2025-08-19 09:31:16.988053+00","2025-08-26 11:14:40.584875+00","2025-09-01 22:05:27.736848+00","2025-09-09 05:38:15.480748+00","2025-09-16 11:08:25.78805+00","2025-09-24 00:59:56.263798+00","2025-09-30 23:02:48.218759+00","2025-10-07 04:23:46.270133+00","2025-10-13 10:14:11.988777+00","2025-10-20 14:06:56.150528+00","2025-10-27 10:54:03.701024+00","2025-11-02 21:57:38.550449+00","2025-11-10 05:40:45.493441+00","2025-11-16 23:22:27.682538+00","2025-11-24 08:46:59.44059+00","2025-11-30 23:18:28.061901+00"}'::text,
	'correlation', '-0.00022358319'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'attname', 'email',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '22'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{user100012@example.com,user102548@example.com,user105478@example.com,user108347@example.com,user111271@example.com,user114260@example.com,user116989@example.com,user119388@example.com,user121744@example.com,user124655@example.com,user127620@example.com,user130339@example.com,user13315@example.com,user135819@example.com,user138495@example.com,user141290@example.com,user143991@example.com,user14661@example.com,user149586@example.com,user152267@example.com,user154812@example.com,user157500@example.com,user160186@example.com,user162906@example.com,user165633@example.com,user168340@example.com,user171001@example.com,user173890@example.com,user176483@example.com,user179188@example.com,user181774@example.com,user184376@example.com,user186988@example.com,user189774@example.com,user192693@example.com,user19531@example.com,user197970@example.com,user200654@example.com,user203756@example.com,user206418@example.com,user209267@example.com,user212138@example.com,user215111@example.com,user217897@example.com,user220512@example.com,user223274@example.com,user226135@example.com,user22881@example.com,user231377@example.com,user234039@example.com,user236684@example.com,user239526@example.com,user242076@example.com,user24495@example.com,user247522@example.com,user250007@example.com,user252814@example.com,user255534@example.com,user258235@example.com,user2610@example.com,user263593@example.com,user266261@example.com,user26903@example.com,user27157@example.com,user27439@example.com,user277233@example.com,user279757@example.com,user282417@example.com,user284950@example.com,user287551@example.com,user290194@example.com,user292637@example.com,user295138@example.com,user297847@example.com,user30515@example.com,user33301@example.com,user35932@example.com,user38507@example.com,user41237@example.com,user4389@example.com,user46698@example.com,user49384@example.com,user52176@example.com,user54993@example.com,user5794@example.com,user60746@example.com,user63274@example.com,user66018@example.com,user68508@example.com,user71173@example.com,user73703@example.com,user76625@example.com,user79211@example.com,user81663@example.com,user84211@example.com,user86996@example.com,user8962@example.com,user9214@example.com,user94745@example.com,user97146@example.com,user99999@example.com}'::text,
	'correlation', '-0.1343658'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'attname', 'id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{7,2876,5605,8452,11248,14301,17404,20529,23534,26435,29280,32339,35341,38280,41301,44225,47333,50400,53594,56824,59992,62823,65916,68691,71680,74698,77682,80568,83454,86462,89374,92346,95182,98192,101091,104115,107525,110492,114016,116955,119533,122245,125543,128812,131882,134916,137784,140885,143896,146899,150153,153105,155923,158901,161966,164957,167989,170912,174046,177155,179895,182691,185642,188580,191828,194740,197759,200757,204303,207496,210469,213586,216784,219737,222640,225930,228899,231812,234829,237811,240755,243796,246782,249433,252677,255677,258712,261782,264744,267723,270564,273661,276819,279705,282593,285472,288350,291249,294186,297163,299991}'::text,
	'correlation', '1'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'attname', 'tenant_id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '4'::integer,
	'n_distinct', '1999'::real,
	'most_common_vals', '{1,593,1106,1839,163,954,586,1255,324,448,544,892,1873}'::text,
	'most_common_freqs', '{0.3957,0.00076666666,0.00073333335,0.00073333335,0.00066666666,0.00066666666,0.00063333334,0.00063333334,0.0006,0.0006,0.0006,0.0006,0.0006}'::real[],
	'histogram_bounds', '{2,20,39,59,80,100,122,140,160,180,199,217,236,255,274,291,311,335,356,376,397,417,439,459,477,499,519,537,561,579,600,618,639,657,676,695,716,735,754,772,792,811,829,849,869,887,908,929,948,967,986,1007,1026,1048,1067,1087,1107,1128,1149,1170,1189,1209,1229,1250,1271,1290,1310,1332,1352,1371,1393,1410,1431,1452,1473,1492,1514,1535,1554,1572,1593,1612,1633,1653,1675,1695,1715,1735,1754,1775,1795,1816,1836,1858,1880,1900,1922,1941,1961,1980,2000}'::text,
	'correlation', '0.16263658'::real
);


--
-- Statistics for Name: order_items; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items',
	'relpages', '29412'::integer,
	'reltuples', '4e+06'::real,
	'relallvisible', '29412'::integer,
	'relallfrozen', '0'::integer
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items',
	'attname', 'order_id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-0.48590675'::real,
	'histogram_bounds', '{15,19515,40578,60886,79478,98316,116781,137564,155925,178555,197469,220932,240687,261739,282260,303235,325287,343245,364428,381747,402229,419991,438386,460469,480276,500340,521140,539074,560658,581419,601912,621304,639671,660577,681900,703011,723920,743865,763483,783299,804356,824661,845929,865709,884388,904450,926560,948635,967109,985708,1005197,1025038,1043552,1062971,1082139,1102143,1121980,1142337,1161873,1182261,1203073,1222628,1241889,1259524,1279207,1297187,1317363,1338591,1359283,1379494,1399447,1417568,1438125,1458914,1480023,1497154,1515350,1536130,1557170,1575993,1598068,1620111,1638767,1659198,1678156,1698580,1719450,1740799,1762130,1783005,1801586,1821187,1840611,1860873,1880194,1898925,1920655,1941196,1961940,1982204,1999970}'::text,
	'correlation', '0.12664323'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items',
	'attname', 'price',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '6'::integer,
	'n_distinct', '9955'::real,
	'most_common_vals', '{20.55,81.44,1.63,26.47,43.19,48.33,67.27,69.13,77.79,79.06,79.25,79.35,80.84,90.19,91.30}'::text,
	'most_common_freqs', '{0.00036666667,0.00036666667,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333,0.00033333333}'::real[],
	'histogram_bounds', '{0.00,1.02,1.95,3.07,4.05,5.04,6.02,7.04,7.96,8.91,9.92,10.86,11.89,12.85,13.73,14.64,15.63,16.70,17.67,18.59,19.49,20.59,21.53,22.47,23.51,24.50,25.45,26.48,27.47,28.55,29.52,30.53,31.47,32.39,33.39,34.40,35.41,36.36,37.36,38.28,39.30,40.34,41.29,42.31,43.38,44.38,45.36,46.43,47.51,48.43,49.36,50.38,51.41,52.47,53.55,54.53,55.61,56.58,57.64,58.67,59.67,60.71,61.70,62.69,63.62,64.52,65.42,66.45,67.42,68.51,69.44,70.46,71.47,72.48,73.45,74.39,75.52,76.55,77.46,78.47,79.59,80.59,81.70,82.86,83.96,85.02,86.11,87.11,88.04,88.97,90.00,91.06,92.06,93.10,94.07,95.02,96.11,97.14,98.07,99.04,100.00}'::text,
	'correlation', '-5.7836492e-05'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items',
	'attname', 'qty',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '4'::integer,
	'n_distinct', '3'::real,
	'most_common_vals', '{2,3,1}'::text,
	'most_common_freqs', '{0.33496666,0.33336666,0.33166668}'::real[],
	'correlation', '0.33553895'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items',
	'attname', 'sku',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '6'::integer,
	'n_distinct', '1998'::real,
	'most_common_vals', '{SKU1,SKU2,SKU3,SKU4,SKU5,SKU6,SKU7,SKU9,SKU8,SKU11,SKU10,SKU13,SKU12,SKU14,SKU16,SKU15,SKU18,SKU17,SKU20,SKU19,SKU24,SKU26,SKU21,SKU23,SKU27,SKU22,SKU29,SKU25,SKU30,SKU34,SKU35,SKU28,SKU38,SKU44,SKU37,SKU41,SKU49,SKU32,SKU45,SKU40,SKU39,SKU31,SKU36,SKU33,SKU46,SKU61,SKU43,SKU47,SKU65,SKU56,SKU51,SKU53,SKU73,SKU50,SKU71,SKU52,SKU54,SKU55,SKU42,SKU74,SKU60,SKU69,SKU57,SKU59,SKU67,SKU95,SKU88,SKU48,SKU99,SKU79,SKU58,SKU84,SKU87,SKU101,SKU102,SKU66,SKU72,SKU94,SKU63,SKU82,SKU85,SKU68,SKU81,SKU117,SKU70,SKU90,SKU93,SKU120,SKU89,SKU62,SKU77,SKU78,SKU111,SKU116,SKU139,SKU75,SKU83,SKU109,SKU126,SKU161}'::text,
	'most_common_freqs', '{0.07946666,0.021233333,0.014833333,0.011533333,0.009633333,0.008933334,0.0074,0.0072333333,0.0066333334,0.0057333335,0.0057,0.005266667,0.0047333334,0.0046,0.0044333334,0.004333333,0.004333333,0.004166667,0.0039333333,0.0038,0.0035333333,0.0035333333,0.0033666666,0.0033,0.0031333333,0.0030333332,0.0029666666,0.0029333334,0.0028666668,0.0026333334,0.0026333334,0.0025,0.0024666667,0.0024333333,0.0024,0.0024,0.0024,0.0023333333,0.0023333333,0.0022666666,0.0022,0.0021666666,0.0021,0.0020666667,0.0020666667,0.0020666667,0.0020333333,0.0020333333,0.0020333333,0.0019666667,0.0019,0.0019,0.0019,0.0018666667,0.0018666667,0.0018333333,0.0018333333,0.0018333333,0.0018,0.0018,0.0017666667,0.0017333333,0.0016666667,0.0016666667,0.0016666667,0.0016666667,0.0016333334,0.0016,0.0016,0.0015666666,0.0015333333,0.0015333333,0.0015333333,0.0015,0.0015,0.0015,0.0015,0.0015,0.0014666667,0.0014666667,0.0014666667,0.0014333334,0.0014333334,0.0014,0.0014,0.0014,0.0014,0.0013666666,0.0013666666,0.0013333333,0.0013333333,0.0013333333,0.0013,0.0013,0.0013,0.0013,0.0013,0.0012666667,0.0012666667,0.0012666667}'::real[],
	'histogram_bounds', '{SKU100,SKU1020,SKU104,SKU1056,SKU1072,SKU1092,SKU1111,SKU113,SKU115,SKU1171,SKU119,SKU1212,SKU1230,SKU1250,SKU1274,SKU1294,SKU1315,SKU1331,SKU1353,SKU1375,SKU140,SKU1420,SKU144,SKU146,SKU1484,SKU1507,SKU1529,SKU155,SKU1573,SKU1595,SKU1620,SKU1642,SKU1665,SKU1689,SKU171,SKU1732,SKU1756,SKU1781,SKU1804,SKU183,SKU1853,SKU188,SKU1905,SKU1932,SKU1957,SKU1981,SKU202,SKU209,SKU217,SKU226,SKU235,SKU245,SKU253,SKU263,SKU273,SKU283,SKU294,SKU304,SKU314,SKU325,SKU337,SKU348,SKU360,SKU372,SKU385,SKU398,SKU409,SKU422,SKU436,SKU449,SKU462,SKU477,SKU493,SKU508,SKU523,SKU537,SKU551,SKU568,SKU585,SKU602,SKU617,SKU634,SKU649,SKU668,SKU687,SKU706,SKU728,SKU747,SKU765,SKU785,SKU800,SKU823,SKU843,SKU861,SKU883,SKU909,SKU922,SKU943,SKU963,SKU98,SKU999}'::text,
	'correlation', '0.0040638475'::real
);


--
-- Statistics for Name: orders; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'relpages', '20764'::integer,
	'reltuples', '2e+06'::real,
	'relallvisible', '20558'::integer,
	'relallfrozen', '0'::integer
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'coupon_code',
	'inherited', 'f'::boolean,
	'null_frac', '0.8999'::real,
	'avg_width', '5'::integer,
	'n_distinct', '9710'::real,
	'histogram_bounds', '{C10,C1082,C1190,C1294,C1372,C1499,C1603,C1700,C1807,C1905,C1995,C2059,C2171,C2254,C2360,C2458,C2522,C2617,C2721,C2802,C2893,C2997,C3115,C3215,C3321,C3415,C3519,C3597,C3700,C3787,C3889,C3942,C4016,C4082,C4174,C4283,C4363,C4431,C452,C4634,C473,C4797,C4869,C4956,C5013,C5117,C5220,C5292,C5390,C5500,C5585,C5709,C5779,C5883,C5932,C6029,C6109,C6187,C6282,C6368,C6455,C6522,C6593,C6677,C6750,C6878,C6961,C7045,C712,C7226,C7332,C7428,C7493,C7607,C7689,C7779,C788,C80,C8089,C8164,C8267,C8348,C8420,C8489,C8572,C8662,C8737,C8802,C89,C8991,C9087,C9179,C9244,C9334,C9404,C9499,C9590,C9690,C9787,C9878,C9998}'::text,
	'correlation', '-0.020031448'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'created_at',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{"2024-01-01 00:44:18.491417+00","2024-01-08 06:02:51.719293+00","2024-01-14 18:19:12.414494+00","2024-01-21 22:34:24.58551+00","2024-01-28 22:24:59.583113+00","2024-02-04 23:38:39.784447+00","2024-02-12 11:18:40.373908+00","2024-02-19 04:03:46.123436+00","2024-02-26 03:19:34.063138+00","2024-03-03 20:36:05.903323+00","2024-03-11 02:31:58.821435+00","2024-03-18 03:50:56.467556+00","2024-03-24 22:17:42.392186+00","2024-04-01 14:27:18.067477+00","2024-04-08 23:59:12.400416+00","2024-04-16 09:23:51.351332+00","2024-04-23 15:23:00.305636+00","2024-04-30 12:19:44.744403+00","2024-05-07 23:20:13.346716+00","2024-05-15 05:20:21.921991+00","2024-05-22 06:23:51.370458+00","2024-05-29 04:11:01.888486+00","2024-06-04 18:44:50.773236+00","2024-06-11 14:55:55.579207+00","2024-06-19 02:29:02.257276+00","2024-06-25 02:11:28.387319+00","2024-07-02 00:00:09.011319+00","2024-07-08 05:20:04.789419+00","2024-07-15 11:13:38.414661+00","2024-07-22 03:13:29.237192+00","2024-07-28 21:11:42.894765+00","2024-08-05 00:46:00.727596+00","2024-08-11 15:06:07.503779+00","2024-08-18 20:32:34.431748+00","2024-08-25 23:53:47.370882+00","2024-09-01 18:08:07.345103+00","2024-09-08 14:37:36.723185+00","2024-09-15 02:38:37.227992+00","2024-09-22 20:13:59.463119+00","2024-09-29 14:37:39.97203+00","2024-10-06 20:19:26.733731+00","2024-10-13 21:46:51.572921+00","2024-10-20 16:25:34.542943+00","2024-10-28 16:52:00.757836+00","2024-11-04 10:03:13.709149+00","2024-11-11 02:08:40.326363+00","2024-11-17 23:53:06.429395+00","2024-11-25 06:21:35.506425+00","2024-12-02 02:14:34.99456+00","2024-12-08 23:55:21.842545+00","2024-12-16 22:44:55.70638+00","2024-12-23 09:49:49.954969+00","2024-12-30 10:10:16.091125+00","2025-01-05 18:09:49.456798+00","2025-01-13 02:58:25.731598+00","2025-01-19 15:53:24.368445+00","2025-01-26 17:35:36.967206+00","2025-02-03 04:46:42.71957+00","2025-02-09 15:47:25.021462+00","2025-02-15 22:35:49.001509+00","2025-02-23 16:28:20.487232+00","2025-03-02 11:20:18.765758+00","2025-03-09 13:31:58.538293+00","2025-03-16 01:25:37.417527+00","2025-03-22 21:51:29.17256+00","2025-03-30 05:09:23.899444+00","2025-04-05 22:50:26.800058+00","2025-04-12 22:43:14.394705+00","2025-04-19 19:43:03.648173+00","2025-04-26 18:36:02.09968+00","2025-05-03 20:14:53.366381+00","2025-05-10 21:35:04.355742+00","2025-05-17 23:15:43.002624+00","2025-05-24 22:54:57.798206+00","2025-05-31 23:28:33.465888+00","2025-06-07 12:23:37.440428+00","2025-06-14 21:01:37.643773+00","2025-06-21 07:01:02.129595+00","2025-06-28 16:47:26.761681+00","2025-07-06 07:14:09.097514+00","2025-07-13 23:55:11.118755+00","2025-07-20 17:27:33.45048+00","2025-07-28 06:01:31.166788+00","2025-08-04 09:31:20.429747+00","2025-08-11 16:49:22.149876+00","2025-08-17 18:18:47.53185+00","2025-08-24 13:25:13.293791+00","2025-08-31 03:14:04.743995+00","2025-09-07 01:43:20.332696+00","2025-09-14 07:06:36.867062+00","2025-09-21 10:59:45.666762+00","2025-09-28 08:47:07.491119+00","2025-10-05 16:24:12.184913+00","2025-10-12 04:33:22.492527+00","2025-10-19 02:07:03.279004+00","2025-10-26 15:00:36.248976+00","2025-11-02 17:30:50.431718+00","2025-11-09 23:13:20.376552+00","2025-11-17 13:35:23.995055+00","2025-11-23 21:02:40.957054+00","2025-11-30 23:36:55.918779+00"}'::text,
	'correlation', '0.004147592'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'customer_id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-0.1243535'::real,
	'histogram_bounds', '{10,2623,5845,8833,11612,14577,17549,20810,23480,26510,29653,32502,35596,38789,41664,44685,48248,51400,54754,57713,60638,63827,66624,69679,72967,75756,78648,81838,85039,88048,91027,94099,96898,99817,102954,106141,109185,112280,115471,118766,121755,124538,127298,130530,133392,136487,139230,142099,145015,147992,151027,154099,157183,159866,162888,165782,169105,172103,175322,178414,181024,184145,187008,190366,193617,196473,199284,201952,204832,207654,210790,213787,216725,219577,222830,225743,228760,231648,234635,237864,240856,243361,246423,249742,253056,255869,258701,261535,264616,267688,270510,273646,276436,279192,282605,285808,288456,291339,294167,297237,299995}'::text,
	'correlation', '0.0021175256'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '8'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{88,19288,39043,58428,77694,99849,119033,139206,158200,178498,200261,219789,240182,260317,279669,300984,320118,340841,361169,382980,403294,423153,443559,465326,484342,504280,524085,543629,562069,581914,601138,621211,642079,662982,681860,704406,725020,744215,765563,784750,805987,825898,844851,863172,883441,903517,921782,941810,963519,983088,1003531,1022072,1040165,1059655,1079082,1099384,1120944,1139676,1160502,1180704,1200813,1220516,1240367,1262378,1283014,1303022,1322148,1341640,1361658,1380954,1400983,1419021,1438954,1458364,1479685,1499551,1517259,1536971,1556646,1576042,1596105,1615983,1636002,1654639,1674549,1698267,1717216,1735787,1756205,1775163,1797783,1819571,1841022,1860382,1879818,1900922,1919217,1939659,1959313,1981201,1999798}'::text,
	'correlation', '0.25869784'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'status',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '9'::integer,
	'n_distinct', '4'::real,
	'most_common_vals', '{delivered,pending,refunded,fraud_review}'::text,
	'most_common_freqs', '{0.9526,0.0424,0.0044333334,0.00056666665}'::real[],
	'correlation', '0.91274965'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'tenant_id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '4'::integer,
	'n_distinct', '2000'::real,
	'most_common_vals', '{1,972,1677,20,1759,1926,108,123,793,1061,124,330,965,1608,1782,1988,1992}'::text,
	'most_common_freqs', '{0.40343332,0.0007,0.0007,0.00066666666,0.00066666666,0.00066666666,0.00063333334,0.00063333334,0.00063333334,0.00063333334,0.0006,0.0006,0.0006,0.0006,0.0006,0.0006,0.0006}'::real[],
	'histogram_bounds', '{2,24,43,62,82,103,122,144,163,182,205,225,245,264,281,304,323,346,367,387,408,427,450,470,489,510,528,549,568,588,608,627,645,666,687,711,729,747,767,788,808,827,846,866,885,904,921,941,961,986,1006,1025,1042,1063,1084,1105,1124,1144,1163,1181,1200,1219,1239,1257,1277,1299,1320,1338,1358,1377,1399,1421,1441,1461,1482,1500,1520,1539,1559,1580,1597,1618,1636,1654,1674,1695,1714,1737,1758,1781,1801,1821,1840,1859,1879,1897,1917,1937,1958,1978,2000}'::text,
	'correlation', '0.16933355'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'attname', 'total',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '6'::integer,
	'n_distinct', '29888'::real,
	'most_common_vals', '{43.37,90.63,98.91}'::text,
	'most_common_freqs', '{0.00023333334,0.00023333334,0.00023333334}'::real[],
	'histogram_bounds', '{0.00,3.08,6.09,9.32,12.56,15.59,18.83,21.59,24.78,27.53,30.44,33.36,36.45,39.70,42.77,45.87,48.73,51.71,54.77,57.40,60.36,63.44,66.28,69.19,72.34,75.60,78.90,81.87,84.65,87.69,90.69,93.72,97.20,100.10,103.23,106.17,109.13,112.01,114.96,118.24,121.31,124.45,127.54,130.61,133.71,136.80,139.79,142.85,145.69,148.91,151.88,154.56,157.66,160.84,163.84,166.57,169.69,172.56,175.40,178.20,181.12,184.07,187.14,190.31,193.40,196.59,199.70,202.54,205.54,208.49,211.62,214.67,217.76,220.52,223.45,226.04,228.90,231.73,234.89,237.91,240.86,243.58,246.62,249.68,252.62,255.99,258.87,261.70,264.60,267.39,270.48,273.35,276.11,279.22,282.39,285.29,288.29,291.27,294.18,297.08,299.99}'::text,
	'correlation', '0.0022203082'::real
);


--
-- Statistics for Name: tenants; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'tenants',
	'relpages', '13'::integer,
	'reltuples', '2000'::real,
	'relallvisible', '13'::integer,
	'relallfrozen', '13'::integer
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'tenants',
	'attname', 'id',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '4'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{1,20,40,60,80,100,120,140,160,180,200,220,240,260,280,300,320,340,360,380,400,420,440,460,480,500,520,540,560,580,600,620,640,660,680,700,720,740,760,780,800,820,840,860,880,900,920,940,960,980,1000,1020,1040,1060,1080,1100,1120,1140,1160,1180,1200,1220,1240,1260,1280,1300,1320,1340,1360,1380,1400,1420,1440,1460,1480,1500,1520,1540,1560,1580,1600,1620,1640,1660,1680,1700,1720,1740,1760,1780,1800,1820,1840,1860,1880,1900,1920,1940,1960,1980,2000}'::text,
	'correlation', '1'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'tenants',
	'attname', 'name',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '10'::integer,
	'n_distinct', '-1'::real,
	'histogram_bounds', '{tenant1,tenant1015,tenant1033,tenant1051,tenant107,tenant1088,tenant1105,tenant1123,tenant1141,tenant116,tenant1178,tenant1196,tenant1213,tenant1231,tenant125,tenant1268,tenant1286,tenant1303,tenant1321,tenant134,tenant1358,tenant1376,tenant1394,tenant1411,tenant143,tenant1448,tenant1466,tenant1484,tenant1501,tenant152,tenant1538,tenant1556,tenant1574,tenant1592,tenant161,tenant1628,tenant1646,tenant1664,tenant1682,tenant170,tenant1718,tenant1736,tenant1754,tenant1772,tenant1790,tenant1808,tenant1826,tenant1844,tenant1862,tenant1880,tenant1899,tenant1916,tenant1934,tenant1952,tenant1970,tenant1989,tenant205,tenant223,tenant241,tenant26,tenant278,tenant296,tenant313,tenant331,tenant35,tenant368,tenant386,tenant403,tenant421,tenant44,tenant458,tenant476,tenant494,tenant511,tenant53,tenant548,tenant566,tenant584,tenant601,tenant62,tenant638,tenant656,tenant674,tenant692,tenant71,tenant728,tenant746,tenant764,tenant782,tenant80,tenant818,tenant836,tenant854,tenant872,tenant890,tenant908,tenant926,tenant944,tenant962,tenant980,tenant999}'::text,
	'correlation', '-0.38968256'::real
);
SELECT * FROM pg_catalog.pg_restore_attribute_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'tenants',
	'attname', 'plan',
	'inherited', 'f'::boolean,
	'null_frac', '0'::real,
	'avg_width', '5'::integer,
	'n_distinct', '2'::real,
	'most_common_vals', '{free,enterprise}'::text,
	'most_common_freqs', '{0.9,0.1}'::real[],
	'correlation', '0.8186477'::real
);


--
-- Statistics for Name: coupons_pkey; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'coupons_pkey',
	'relpages', '94'::integer,
	'reltuples', '20000'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: customers_email_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers_email_idx',
	'relpages', '2610'::integer,
	'reltuples', '300000'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: customers_pkey; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers_pkey',
	'relpages', '825'::integer,
	'reltuples', '300000'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: customers_tenant_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers_tenant_idx',
	'relpages', '322'::integer,
	'reltuples', '300000'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: order_items_order_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items_order_idx',
	'relpages', '11700'::integer,
	'reltuples', '4e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: order_items_sku_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'order_items_sku_idx',
	'relpages', '3959'::integer,
	'reltuples', '4e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: orders_customer_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders_customer_idx',
	'relpages', '3415'::integer,
	'reltuples', '2e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: orders_pkey; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders_pkey',
	'relpages', '7126'::integer,
	'reltuples', '2e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: orders_status_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders_status_idx',
	'relpages', '1588'::integer,
	'reltuples', '2e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: orders_tenant_created_idx; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders_tenant_created_idx',
	'relpages', '9876'::integer,
	'reltuples', '2e+06'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- Statistics for Name: tenants_pkey; Type: STATISTICS DATA; Schema: public; Owner: -
--

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'tenants_pkey',
	'relpages', '8'::integer,
	'reltuples', '2000'::real,
	'relallvisible', '0'::integer,
	'relallfrozen', '0'::integer
);


--
-- PostgreSQL database dump complete
--

\unrestrict stuntdb

