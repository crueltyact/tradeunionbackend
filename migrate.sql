ALTER TABLE config.classification
    ALTER COLUMN domain TYPE text;

ALTER TABLE config.rule
    ALTER COLUMN domain TYPE text;

DROP TYPE config.slo_domain;

INSERT INTO config.domain (id, domain_name, enabled, namespace, version) VALUES (2, 'crypto-core', true, 'aifory-pay', 1);
INSERT INTO config.domain (id, domain_name, enabled, namespace, version) VALUES (3, 'neo-bank', true, 'aifory-pay', 1);
INSERT INTO config.domain (id, domain_name, enabled, namespace, version) VALUES (4, 'auto-exchange', true, 'aifory-pay', 1);
INSERT INTO config.domain (id, domain_name, enabled, namespace, version) VALUES (5, 'platform', true, 'aifory-pay', 1);
INSERT INTO config.domain (id, domain_name, enabled, namespace, version) VALUES (1, 'nal', true, 'aifory-pay', 1);


UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'api', class = 'low', pattern = 'api:platform:low:*', match_type = 'regexp' WHERE pattern = 'api:low:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'api', class = 'critical', pattern = 'api:platform:critical:*', match_type = 'regexp' WHERE pattern = 'api:critical:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'api', class = 'medium', pattern = 'api:platform:medium:*', match_type = 'regexp' WHERE pattern = 'api:medium:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'system', class = 'medium', pattern = 'system:platform:medium:*', match_type = 'regexp' WHERE pattern = 'system:medium:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'system', class = 'critical', pattern = 'system:platform:critical:*', match_type = 'regexp' WHERE pattern = 'system:critical:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'app', class = 'low', pattern = 'app:platofrm:low:*', match_type = 'regexp' WHERE pattern = 'app:low:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'app', class = 'critical', pattern = 'app:platofrm:critical:*', match_type = 'regexp' WHERE pattern = 'app:critical:*';
UPDATE config.classification SET pipeline_id = 2, domain = 'platform', app = 'app', class = 'medium', pattern = 'app:platofrm:medium:*', match_type = 'regexp' WHERE pattern = 'app:medium:*';


UPDATE config.prom_query SET type = 'simple',  interval = '1m', "offset" = '24h', label = 'kafka_proccess_running', providervalue = 'system:platform:critical', executedelay = '30s', datasource = 'prometheus' WHERE query = 'sum(processes_running{service="ap-kafka", job="do-dbaas-kafka"})';
UPDATE config.prom_query SET type = 'simple', interval = '1m', "offset" = '24h', label = 'api_server_up', providervalue = 'system:platform:critical', executedelay = '30s', datasource = 'victoriametrics' WHERE query = 'sum(up{job=~"kubernetes-apiservers|apiserver"})';
UPDATE config.prom_query SET type = 'simple', interval = '1m', "offset" = '24h', label = 'fail_cloudfare_rate', providervalue = 'system:platform:critical', executedelay = '30s', datasource = 'victoriametrics' WHERE query = '100 - sum(increase(cloudflare_zone_requests_status{status=~"5.."}[1m])) / sum(increase(cloudflare_zone_requests_status[1m])) * 100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'network_receive_errors_total', providervalue = 'system:platform:critical', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum(increase(container_network_receive_errors_total{namespace="aifory-pay", cluster="do-prod-ap"}[5m]))';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'coredns_health_status', providervalue = 'system:platform:critical', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum (up{job=~"coredns", cluster=~"do-prod-ap"})';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'apiserver_request_total', providervalue = 'system:platform:critical', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100-(sum (rate(apiserver_request_total{cluster=~"do-prod-ap", code=~"5.."}[5m])) / sum (rate(apiserver_request_total{cluster=~"do-prod-ap"}[5m]))*100)';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'http_fail_client_rate_400', providervalue = 'api:platform:medium:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = e'100 - (sum(increase(slo_exporter_http_requests_total{status_code=~"^4..", path=~".*lk.*"}[1m])) or 0)
/
sum(increase(slo_exporter_http_requests_total{path=~".*lk.*"}[1m]))
* 100
';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'http_fail_admin_rate', providervalue = 'api:platform:low:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100 - (sum(increase(slo_exporter_http_requests_total{status_code=~"^5..", path=~".*admin.*"}[1m])) or 0)/sum(increase(slo_exporter_http_requests_total{path=~".*admin.*"}[1m]))* 100';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'http_fail_client_rate', providervalue = 'api:platform:critical:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = e'100 - (sum(increase(slo_exporter_http_requests_total{status_code=~"^5..", path=~".*lk.*"}[1m])) or 0)
/
sum(increase(slo_exporter_http_requests_total{path=~".*lk.*"}[1m]))
* 100
';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'coredns_servfail_rate', providervalue = 'system:platform:critical:coredns', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum(rate(coredns_dns_responses_total{rcode!="SERVFAIL"}[1m]))/ sum(rate(coredns_dns_responses_total{}[1m])) * 100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'http_time_admin_execution', providervalue = 'api:platform:low:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'rate(sum by (path) (slo_exporter_response_time_second_sum{path=~".*admin.*"}[1m])) ';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'http_time_client_execution', providervalue = 'api:platform:medium:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'rate(sum by (path) (slo_exporter_response_time_second_sum{path=~".*lk.*"}[1m])) ';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'http_fail_admin_rate_400', providervalue = 'api:platform:low:http', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100 - (sum(increase(slo_exporter_http_requests_total{status_code=~"^4..", path=~".*admin.*"}[1m])) or 0)/sum(increase(slo_exporter_http_requests_total{path=~".*admin.*"}[1m]))* 100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'kafka_consumer_group_lag', providervalue = 'system:platform:medium:kafka', executedelay = '0s', datasource = 'prometheus' WHERE query = 'sum(kafka_consumer_group_rep_lag{service="ap-kafka", job="do-dbaas-kafka", topic!="ap_outgoing_chat_events"})';
UPDATE config.prom_query SET type = 'simple',  interval = '10s', "offset" = '0h', label = 'fail_events_rate_overall', providervalue = 'system:platform:medium:core', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '                       avg(                         (                             increase(slo_domain_slo_class_slo_app_event_key:slo_events_total{result="fail"}[1m]) == bool 0                         ) * 1 or (                             increase(slo_domain_slo_class_slo_app_event_key:slo_events_total{result="fail"}[1m]) > bool 0                         ) * 0 or vector(1)                       ) == bool 1 or 0              ';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'fail_events_rate', providervalue = 'system:platform:medium:core', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '                       avg(                         (                             increase(slo_domain_slo_class_slo_app_event_key:slo_events_total{result="fail", slo_class!="low"}[1m]) == bool 0                         ) * 1 or (                             increase(slo_domain_slo_class_slo_app_event_key:slo_events_total{result="fail", slo_class!="low"}[1m]) > bool 0                         ) * 0 or vector(1)                       ) == bool 1 or 0              ';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'antifraud_ban_errors', providervalue = 'app:platform:critical:antifraud', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100 - ((sum by (fraud_producer) (increase(antifraud_ban_errors_total[1d]))/sum by (fraud_producer) (increase(antifraud_bans_total[1d])) * 100) or 0)';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'tg_messages_failure_rate_client', providervalue = 'app:platform:medium:telegram', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100-sum by (namespace, destination) (rate(tg_result_messages_total{destination="client", result="fail"}[1d]))/sum by (namespace, destination) (rate(tg_result_messages_total{destination="client"}[1d]))*100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'tg_messages_failure_rate_admin', providervalue = 'app:platform:low:telegram', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100-sum by (namespace, destination) (rate(tg_result_messages_total{destination="admin", result="fail"}[1d]))/sum by (namespace, destination) (rate(tg_result_messages_total{destination="admin"}[1d]))*100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'auth_email_send', providervalue = 'app:platform:critical', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum(increase(auth_email_send_result_total{result="success"}[5m])) / sum(increase(auth_email_send_result_total[5m])) * 100';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'fail_batch_rate', providervalue = 'app:platform:critical:abs', executedelay = '0s', datasource = 'victoriametrics' WHERE query = '100 - ((sum(increase(abs_entries_creation_errors_total[1m]) or 0) / sum(increase(abs_entries_create_attempts_total[1m])) * 100) or 0)';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'fail_overtime_batch_execution', providervalue = 'app:platform:critical:abs', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum by (le) (rate(abs_batch_execution_time_duration_bucket[1m]))';
UPDATE config.prom_query SET type = 'simple', interval = '10s', "offset" = '0h', label = 'outbox_lag', providervalue = 'app:platform:low:outbox', executedelay = '0s', datasource = 'victoriametrics' WHERE query = 'sum(outbox_created_events_count)- sum(outbox_sent_events_count) ';
