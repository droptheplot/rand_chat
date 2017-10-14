CREATE OR REPLACE VIEW charts AS
SELECT (now() - (n || 'days')::interval)::date, COUNT(messages.id)
FROM generate_series(0,30) AS n
LEFT OUTER JOIN messages ON messages.created_at > (now() - (n || 'days')::interval)::date
                        AND messages.created_at <= (now() - (n - 1 || 'days')::interval)::date
GROUP BY date
ORDER BY date ASC;
