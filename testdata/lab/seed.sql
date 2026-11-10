-- ported from regresql's poc/regressions-lab/db/gen.sql
select setseed(0.42);

\set ntenants 2000
\set ncust 300000
\set ncoupons 20000
\set norders 2000000
\set nsku 2000

insert into tenants select i, 'tenant'||i, case when i%10=0 then 'enterprise' else 'free' end from generate_series(1,:ntenants) i;
-- tenant 1 owns ~40% of customers
insert into customers select i, case when random()<0.4 then 1 else 2+floor(random()*(:ntenants-1))::int end, 'user'||i||'@example.com', timestamptz '2024-01-01' + random()*interval '700 days' from generate_series(1,:ncust) i;
insert into coupons select 'C'||i, 5+i%20 from generate_series(1,:ncoupons) i;
insert into orders select i, c.tenant_id, c.id,
  case when r<0.0005 then 'fraud_review' when r<0.005 then 'refunded' when r<0.045 then 'pending' else 'delivered' end,
  round((random()*300)::numeric,2), timestamptz '2024-01-01' + random()*interval '700 days',
  case when random()<0.1 then 'C'||(1+floor(random()*(:ncoupons/2))::int) end
from (select i, random() r, 1+floor(random()*:ncust)::bigint cid from generate_series(1,:norders) i) g join customers c on c.id=g.cid;
insert into order_items select o.id, 'SKU'||(1+floor(power(random(),3)*:nsku)::int), 1+floor(random()*3)::int, round((random()*100)::numeric,2) from orders o, generate_series(1,2);
analyze;
