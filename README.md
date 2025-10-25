# traffic_guard

```
traffic_guard
├─ deploy
│  ├─ docker-compose
│  │  └─ docker-compose.yml
│  ├─ k8s
│  │  └─ README.md
│  ├─ Makefile
│  └─ README.md
├─ docs
│  ├─ categories.md
│  ├─ proto-via-kafka.md
│  └─ README.md
├─ Makefile
├─ README.md
└─ src
   ├─ go
   │  ├─ api-gw
   │  │  ├─ .DS_Store
   │  │  ├─ build
   │  │  │  ├─ docker-compose.local.yml
   │  │  │  └─ Dockerfile
   │  │  ├─ cmd
   │  │  │  ├─ api-gw
   │  │  │  │  └─ main.go
   │  │  │  └─ server
   │  │  │     └─ main.go
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yml
   │  │  ├─ data
   │  │  │  └─ keycloak
   │  │  │     └─ key.RS256.pub.pem
   │  │  ├─ deploy
   │  │  │  ├─ docker-compose.yml
   │  │  │  └─ Dockerfile
   │  │  ├─ docs
   │  │  │  ├─ docs.go
   │  │  │  ├─ swagger.json
   │  │  │  └─ swagger.yaml
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ apigw
   │  │  │  │     └─ apigw.go
   │  │  │  └─ controllers
   │  │  │     └─ http
   │  │  │        └─ v1
   │  │  │           ├─ blog_handles.go
   │  │  │           ├─ ctxcontrol_handles.go
   │  │  │           ├─ err_interface.go
   │  │  │           ├─ healthcheck.go
   │  │  │           ├─ metadata
   │  │  │           │  └─ usermeta.go
   │  │  │           ├─ middleware
   │  │  │           │  ├─ auth_header.go
   │  │  │           │  ├─ cors.go
   │  │  │           │  ├─ logging.go
   │  │  │           │  ├─ request_id.go
   │  │  │           │  └─ token.go
   │  │  │           ├─ router.go
   │  │  │           ├─ server.go
   │  │  │           ├─ usercontrol_handles.go
   │  │  │           ├─ utils
   │  │  │           │  ├─ authinfo.go
   │  │  │           │  ├─ request_id.go
   │  │  │           │  └─ user_meta.go
   │  │  │           └─ values
   │  │  │              └─ values.go
   │  │  ├─ Makefile
   │  │  ├─ pkg
   │  │  │  ├─ blogserv
   │  │  │  │  ├─ blogserv_client.go
   │  │  │  │  └─ operations
   │  │  │  │     ├─ get_api_v1_healthcheck_parameters.go
   │  │  │  │     ├─ get_api_v1_healthcheck_responses.go
   │  │  │  │     ├─ get_api_v1_logs_parameters.go
   │  │  │  │     ├─ get_api_v1_logs_responses.go
   │  │  │  │     └─ operations_client.go
   │  │  │  ├─ ctxcontrol
   │  │  │  │  ├─ contexts
   │  │  │  │  │  ├─ contexts_client.go
   │  │  │  │  │  ├─ delete_v1_contexts_context_id_parameters.go
   │  │  │  │  │  ├─ delete_v1_contexts_context_id_responses.go
   │  │  │  │  │  ├─ get_v1_contexts_context_id_parameters.go
   │  │  │  │  │  ├─ get_v1_contexts_context_id_responses.go
   │  │  │  │  │  ├─ get_v1_contexts_count_parameters.go
   │  │  │  │  │  ├─ get_v1_contexts_count_responses.go
   │  │  │  │  │  ├─ get_v1_contexts_free_ports_parameters.go
   │  │  │  │  │  ├─ get_v1_contexts_free_ports_responses.go
   │  │  │  │  │  ├─ get_v1_contexts_parameters.go
   │  │  │  │  │  ├─ get_v1_contexts_responses.go
   │  │  │  │  │  ├─ post_v1_contexts_parameters.go
   │  │  │  │  │  ├─ post_v1_contexts_responses.go
   │  │  │  │  │  ├─ put_v1_contexts_context_id_parameters.go
   │  │  │  │  │  └─ put_v1_contexts_context_id_responses.go
   │  │  │  │  └─ ctxcontrol_client.go
   │  │  │  ├─ models
   │  │  │  │  ├─ dto_api_error.go
   │  │  │  │  ├─ dto_auth_webhook_request.go
   │  │  │  │  ├─ dto_context_list_response.go
   │  │  │  │  ├─ dto_count_response.go
   │  │  │  │  ├─ dto_count_user_by_role_response.go
   │  │  │  │  ├─ dto_count_user_for_context_response.go
   │  │  │  │  ├─ dto_create_context_request.go
   │  │  │  │  ├─ dto_error_response.go
   │  │  │  │  ├─ dto_log_user.go
   │  │  │  │  ├─ dto_pagination_meta.go
   │  │  │  │  ├─ dto_roles_list_response.go
   │  │  │  │  ├─ dto_success_response.go
   │  │  │  │  ├─ dto_update_context_request.go
   │  │  │  │  ├─ dto_users_list_response.go
   │  │  │  │  ├─ dto_user_create_request.go
   │  │  │  │  ├─ dto_user_reset_pass.go
   │  │  │  │  ├─ dto_user_update_data.go
   │  │  │  │  ├─ dto_user_update_pass.go
   │  │  │  │  ├─ dto_vectro_free_ports_response.go
   │  │  │  │  ├─ models_api_error.go
   │  │  │  │  ├─ models_business_log.go
   │  │  │  │  ├─ models_context.go
   │  │  │  │  ├─ models_data_limits.go
   │  │  │  │  ├─ models_dto_success_response.go
   │  │  │  │  ├─ models_json_b.go
   │  │  │  │  ├─ models_logs.go
   │  │  │  │  ├─ models_meta.go
   │  │  │  │  ├─ models_role.go
   │  │  │  │  ├─ models_token.go
   │  │  │  │  └─ models_user.go
   │  │  │  ├─ slogger
   │  │  │  │  ├─ error.go
   │  │  │  │  ├─ logctx.go
   │  │  │  │  ├─ README.md
   │  │  │  │  └─ slogger.go
   │  │  │  └─ usercontrol
   │  │  │     ├─ auth
   │  │  │     │  ├─ auth_client.go
   │  │  │     │  ├─ get_v1_auth_sign_out_parameters.go
   │  │  │     │  ├─ get_v1_auth_sign_out_responses.go
   │  │  │     │  ├─ post_v1_auth_refresh_parameters.go
   │  │  │     │  ├─ post_v1_auth_refresh_responses.go
   │  │  │     │  ├─ post_v1_auth_sign_in_parameters.go
   │  │  │     │  ├─ post_v1_auth_sign_in_responses.go
   │  │  │     │  ├─ post_v1_auth_sign_out_parameters.go
   │  │  │     │  └─ post_v1_auth_sign_out_responses.go
   │  │  │     ├─ roles
   │  │  │     │  ├─ delete_v1_roles_context_parameters.go
   │  │  │     │  ├─ delete_v1_roles_context_responses.go
   │  │  │     │  ├─ get_v1_roles_count_parameters.go
   │  │  │     │  ├─ get_v1_roles_count_responses.go
   │  │  │     │  ├─ get_v1_roles_parameters.go
   │  │  │     │  ├─ get_v1_roles_responses.go
   │  │  │     │  ├─ put_v1_roles_context_parameters.go
   │  │  │     │  ├─ put_v1_roles_context_responses.go
   │  │  │     │  └─ roles_client.go
   │  │  │     ├─ usercontrol_client.go
   │  │  │     └─ users
   │  │  │        ├─ delete_v1_users_user_id_parameters.go
   │  │  │        ├─ delete_v1_users_user_id_responses.go
   │  │  │        ├─ delete_v1_users_user_id_roles_parameters.go
   │  │  │        ├─ delete_v1_users_user_id_roles_responses.go
   │  │  │        ├─ get_v1_users_count_by_role_parameters.go
   │  │  │        ├─ get_v1_users_count_by_role_responses.go
   │  │  │        ├─ get_v1_users_count_context_context_id_parameters.go
   │  │  │        ├─ get_v1_users_count_context_context_id_responses.go
   │  │  │        ├─ get_v1_users_count_parameters.go
   │  │  │        ├─ get_v1_users_count_responses.go
   │  │  │        ├─ get_v1_users_parameters.go
   │  │  │        ├─ get_v1_users_profile_parameters.go
   │  │  │        ├─ get_v1_users_profile_responses.go
   │  │  │        ├─ get_v1_users_responses.go
   │  │  │        ├─ get_v1_users_user_id_parameters.go
   │  │  │        ├─ get_v1_users_user_id_responses.go
   │  │  │        ├─ post_v1_users_parameters.go
   │  │  │        ├─ post_v1_users_responses.go
   │  │  │        ├─ put_v1_users_profile_pass_parameters.go
   │  │  │        ├─ put_v1_users_profile_pass_responses.go
   │  │  │        ├─ put_v1_users_user_id_parameters.go
   │  │  │        ├─ put_v1_users_user_id_pass_otp_parameters.go
   │  │  │        ├─ put_v1_users_user_id_pass_otp_responses.go
   │  │  │        ├─ put_v1_users_user_id_responses.go
   │  │  │        ├─ put_v1_users_user_id_roles_parameters.go
   │  │  │        ├─ put_v1_users_user_id_roles_responses.go
   │  │  │        └─ users_client.go
   │  │  └─ README.md
   │  ├─ blog
   │  │  ├─ .DS_Store
   │  │  ├─ build
   │  │  │  ├─ docker-compose.yaml
   │  │  │  └─ Dockerfile
   │  │  ├─ cmd
   │  │  │  └─ main.go
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yml
   │  │  ├─ docs
   │  │  │  ├─ docs.go
   │  │  │  ├─ swagger.json
   │  │  │  └─ swagger.yaml
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ fiermon-blog.go
   │  │  │  ├─ models
   │  │  │  │  ├─ business_logs.go
   │  │  │  │  ├─ dto_Sucsess_Response.go
   │  │  │  │  ├─ response.go
   │  │  │  │  └─ usermeta.go
   │  │  │  ├─ repository
   │  │  │  │  ├─ postgresGorm
   │  │  │  │  │  ├─ initPG.go
   │  │  │  │  │  └─ read.go
   │  │  │  │  └─ repository_test.go
   │  │  │  ├─ server
   │  │  │  │  ├─ dto
   │  │  │  │  │  ├─ responses.go
   │  │  │  │  │  └─ validate.go
   │  │  │  │  ├─ handlers.go
   │  │  │  │  ├─ middleware
   │  │  │  │  │  ├─ auth.go
   │  │  │  │  │  ├─ auth_header.go
   │  │  │  │  │  ├─ cors.go
   │  │  │  │  │  ├─ logging.go
   │  │  │  │  │  ├─ recovery.go
   │  │  │  │  │  └─ request_id.go
   │  │  │  │  ├─ middleware.go
   │  │  │  │  ├─ router.go
   │  │  │  │  ├─ server.go
   │  │  │  │  ├─ utils
   │  │  │  │  │  ├─ authinfo.go
   │  │  │  │  │  ├─ request_id.go
   │  │  │  │  │  └─ user_meta.go
   │  │  │  │  └─ values
   │  │  │  │     └─ values.go
   │  │  │  └─ service
   │  │  │     ├─ get.go
   │  │  │     └─ service.go
   │  │  ├─ Makefile
   │  │  ├─ pkg
   │  │  │  └─ fslog
   │  │  │     ├─ error.go
   │  │  │     ├─ fslogger.go
   │  │  │     ├─ logctx.go
   │  │  │     ├─ README.md
   │  │  │     └─ wsl
   │  │  │        └─ wraper.go
   │  │  └─ README.md
   │  ├─ dashboard_serv
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yml
   │  │  ├─ docs
   │  │  │  └─ docs.json
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ dashboard_serv
   │  │  │  │     └─ dashboard_serv.go
   │  │  │  ├─ controllers
   │  │  │  │  └─ http
   │  │  │  │     └─ v1
   │  │  │  │        ├─ dto
   │  │  │  │        │  ├─ meta.go
   │  │  │  │        │  ├─ requests.go
   │  │  │  │        │  └─ responses.go
   │  │  │  │        ├─ handlers.go
   │  │  │  │        ├─ router.go
   │  │  │  │        └─ server.go
   │  │  │  └─ usecase
   │  │  │     └─ usecase.go
   │  │  └─ Makefile
   │  ├─ Makefile
   │  ├─ tg-analytics
   │  │  ├─ build
   │  │  │  ├─ config.yaml
   │  │  │  ├─ docker-compose.yaml
   │  │  │  └─ Dockerfile
   │  │  ├─ cmd
   │  │  │  └─ analytics
   │  │  │     └─ main.go
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yaml
   │  │  ├─ deploy
   │  │  │  └─ local
   │  │  │     └─ docker-compose.yaml
   │  │  ├─ docs
   │  │  │  ├─ docs.go
   │  │  │  ├─ swagger.json
   │  │  │  └─ swagger.yaml
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ analytics
   │  │  │  │     └─ analytics.go
   │  │  │  ├─ controllers
   │  │  │  │  └─ http
   │  │  │  │     └─ v1
   │  │  │  │        ├─ dto
   │  │  │  │        │  ├─ requests.go
   │  │  │  │        │  ├─ responses.go
   │  │  │  │        │  └─ session.go
   │  │  │  │        ├─ handlers.go
   │  │  │  │        ├─ healthcheck.go
   │  │  │  │        ├─ middleware
   │  │  │  │        │  ├─ cors.go
   │  │  │  │        │  └─ validation.go
   │  │  │  │        ├─ router.go
   │  │  │  │        ├─ server.go
   │  │  │  │        └─ usecase_interface.go
   │  │  │  ├─ models
   │  │  │  │  ├─ filters.go
   │  │  │  │  ├─ gorm_models.go
   │  │  │  │  ├─ metrics.go
   │  │  │  │  └─ models.go
   │  │  │  ├─ pkg
   │  │  │  │  └─ status
   │  │  │  │     └─ status.go
   │  │  │  ├─ repo
   │  │  │  │  └─ postgresql
   │  │  │  │     ├─ common.go
   │  │  │  │     ├─ dashboards.go
   │  │  │  │     ├─ detections.go
   │  │  │  │     ├─ metrics.go
   │  │  │  │     ├─ postgresql.go
   │  │  │  │     └─ sessions.go
   │  │  │  └─ usecase
   │  │  │     ├─ common.go
   │  │  │     ├─ dashboard.go
   │  │  │     ├─ dashboards.go
   │  │  │     ├─ detections.go
   │  │  │     ├─ pg_interface.go
   │  │  │     ├─ sessions.go
   │  │  │     └─ usecase.go
   │  │  ├─ Makefile
   │  │  ├─ pkg
   │  │  │  ├─ pgorm
   │  │  │  │  └─ pgorm
   │  │  │  │     ├─ options.go
   │  │  │  │     ├─ postgres.go
   │  │  │  │     └─ README.md
   │  │  │  ├─ slogger
   │  │  │  │  ├─ error.go
   │  │  │  │  ├─ logctx.go
   │  │  │  │  ├─ README.md
   │  │  │  │  ├─ slogger.go
   │  │  │  │  └─ wsl
   │  │  │  │     └─ wraper.go
   │  │  │  └─ trparser
   │  │  │     ├─ trparser.go
   │  │  │     └─ t_range.go
   │  │  └─ README.md
   │  ├─ tg-etl
   │  │  ├─ build
   │  │  │  ├─ config.yaml
   │  │  │  ├─ docker-compose.yaml
   │  │  │  └─ Dockerfile
   │  │  ├─ cmd
   │  │  │  └─ etl
   │  │  │     └─ main.go
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yaml
   │  │  ├─ deploy
   │  │  │  ├─ local
   │  │  │  │  └─ docker-compose.yaml
   │  │  │  └─ release
   │  │  ├─ doc
   │  │  │  ├─ db.dbml
   │  │  │  ├─ db_structure.md
   │  │  │  ├─ image-1.png
   │  │  │  ├─ image-2.png
   │  │  │  ├─ image.png
   │  │  │  └─ ksu_db_structure.dbml
   │  │  ├─ docs
   │  │  │  ├─ docs.go
   │  │  │  ├─ swagger.json
   │  │  │  └─ swagger.yaml
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ etl
   │  │  │  │     └─ etl.go
   │  │  │  ├─ controllers
   │  │  │  │  ├─ etl
   │  │  │  │  │  └─ etl.go
   │  │  │  │  └─ http
   │  │  │  │     └─ v1
   │  │  │  │        ├─ dto
   │  │  │  │        │  └─ responses.go
   │  │  │  │        ├─ handlers.go
   │  │  │  │        ├─ middleware
   │  │  │  │        │  └─ cors.go
   │  │  │  │        ├─ router.go
   │  │  │  │        └─ server.go
   │  │  │  ├─ models
   │  │  │  │  ├─ etl_models.go
   │  │  │  │  ├─ kafka.go
   │  │  │  │  ├─ ksu_models.go
   │  │  │  │  └─ lists.go
   │  │  │  ├─ repo
   │  │  │  │  ├─ category
   │  │  │  │  │  └─ category.go
   │  │  │  │  ├─ kafka
   │  │  │  │  │  ├─ client.go
   │  │  │  │  │  └─ interface.go
   │  │  │  │  ├─ m_cache
   │  │  │  │  │  └─ m_cache.go
   │  │  │  │  ├─ parser
   │  │  │  │  │  └─ parser.go
   │  │  │  │  └─ postgresql
   │  │  │  │     ├─ elt_repo.go
   │  │  │  │     ├─ ksu_repo.go
   │  │  │  │     └─ repo_pg.go
   │  │  │  ├─ usecase
   │  │  │  │  ├─ handlers
   │  │  │  │  ├─ interface.go
   │  │  │  │  ├─ kafka.go
   │  │  │  │  ├─ kafka_interface.go
   │  │  │  │  ├─ processor.go
   │  │  │  │  ├─ query
   │  │  │  │  │  ├─ category.go
   │  │  │  │  │  ├─ device.go
   │  │  │  │  │  ├─ domain.go
   │  │  │  │  │  ├─ ksu.go
   │  │  │  │  │  ├─ list.go
   │  │  │  │  │  ├─ pg_interface.go
   │  │  │  │  │  ├─ query_usecase.go
   │  │  │  │  │  ├─ session.go
   │  │  │  │  │  ├─ source.go
   │  │  │  │  │  └─ url.go
   │  │  │  │  └─ usecase.go
   │  │  │  └─ utils
   │  │  │     └─ utils.go
   │  │  ├─ Makefile
   │  │  ├─ pkg
   │  │  │  ├─ cslogger
   │  │  │  │  └─ cslogger.go
   │  │  │  ├─ db_structurer
   │  │  │  │  └─ db_structurer.go
   │  │  │  ├─ pgorm
   │  │  │  │  ├─ options.go
   │  │  │  │  └─ postgres.go
   │  │  │  └─ slogger
   │  │  │     ├─ error.go
   │  │  │     ├─ logctx.go
   │  │  │     ├─ README.md
   │  │  │     ├─ slogger.go
   │  │  │     └─ wsl
   │  │  │        └─ wraper.go
   │  │  ├─ README.md
   │  │  └─ test_data
   │  │     └─ security_log-2025-9-20_15-58.csv
   │  ├─ userctrl
   │  │  ├─ bizlog.json
   │  │  ├─ build
   │  │  │  ├─ docker-compose.local.yml
   │  │  │  └─ Dockerfile
   │  │  ├─ cmd
   │  │  │  └─ usercontrol
   │  │  │     └─ main.go
   │  │  ├─ config
   │  │  │  ├─ config.go
   │  │  │  └─ config.yml
   │  │  ├─ deploy
   │  │  │  ├─ docker-compose.local.yml
   │  │  │  └─ Dockerfile
   │  │  ├─ doc
   │  │  │  ├─ keycloakclient
   │  │  │  │  ├─ attribute.png
   │  │  │  │  ├─ client_role.png
   │  │  │  │  └─ keycloakclient.md
   │  │  │  └─ swagger.yaml
   │  │  ├─ docs
   │  │  │  ├─ docs.go
   │  │  │  ├─ swagger.json
   │  │  │  └─ swagger.yaml
   │  │  ├─ go.mod
   │  │  ├─ go.sum
   │  │  ├─ internal
   │  │  │  ├─ app
   │  │  │  │  └─ usercontrol
   │  │  │  │     └─ usercontrol.go
   │  │  │  ├─ controllers
   │  │  │  │  └─ http
   │  │  │  │     └─ v1
   │  │  │  │        ├─ auth_handlers.go
   │  │  │  │        ├─ dto
   │  │  │  │        │  ├─ requests.go
   │  │  │  │        │  └─ responses.go
   │  │  │  │        ├─ errors.go
   │  │  │  │        ├─ healthcheck.go
   │  │  │  │        ├─ middleware
   │  │  │  │        │  ├─ auth.go
   │  │  │  │        │  ├─ auth_header.go
   │  │  │  │        │  ├─ cors.go
   │  │  │  │        │  ├─ logging.go
   │  │  │  │        │  ├─ recovery.go
   │  │  │  │        │  └─ request_id.go
   │  │  │  │        ├─ roles.go
   │  │  │  │        ├─ router.go
   │  │  │  │        ├─ server.go
   │  │  │  │        ├─ usecase_interface.go
   │  │  │  │        ├─ user_handlers.go
   │  │  │  │        ├─ utils
   │  │  │  │        │  ├─ authinfo.go
   │  │  │  │        │  ├─ request_id.go
   │  │  │  │        │  └─ user_meta.go
   │  │  │  │        └─ values
   │  │  │  │           └─ values.go
   │  │  │  ├─ models
   │  │  │  │  ├─ errors.go
   │  │  │  │  └─ models.go
   │  │  │  ├─ repo
   │  │  │  │  └─ blogrepo
   │  │  │  │     └─ postgresql.go
   │  │  │  └─ usecase
   │  │  │     ├─ auth.go
   │  │  │     ├─ blog.go
   │  │  │     ├─ roles.go
   │  │  │     ├─ usecase.go
   │  │  │     ├─ users.go
   │  │  │     └─ utils
   │  │  │        └─ utils.go
   │  │  ├─ Makefile
   │  │  ├─ pkg
   │  │  │  ├─ bizlogger
   │  │  │  │  ├─ bizlog.go
   │  │  │  │  ├─ file.go
   │  │  │  │  ├─ logger.go
   │  │  │  │  ├─ queue.go
   │  │  │  │  └─ README.md
   │  │  │  ├─ grafanaclient
   │  │  │  │  ├─ grafanaclient.go
   │  │  │  │  ├─ org.go
   │  │  │  │  ├─ users.go
   │  │  │  │  └─ validate.go
   │  │  │  ├─ grafcookier
   │  │  │  │  ├─ grafcookier.go
   │  │  │  │  └─ grafcookier_test.go
   │  │  │  ├─ keycloakclient
   │  │  │  │  ├─ auth.go
   │  │  │  │  ├─ errors.go
   │  │  │  │  ├─ groups.go
   │  │  │  │  ├─ groups_test.go
   │  │  │  │  ├─ keycloakclient.go
   │  │  │  │  ├─ keycloakclient_test.go
   │  │  │  │  ├─ roles.go
   │  │  │  │  ├─ roles_test.go
   │  │  │  │  ├─ users.go
   │  │  │  │  ├─ users_test.go
   │  │  │  │  └─ validate.go
   │  │  │  ├─ pgorm
   │  │  │  │  ├─ options.go
   │  │  │  │  ├─ postgres.go
   │  │  │  │  └─ README.md
   │  │  │  └─ slogger
   │  │  │     ├─ error.go
   │  │  │     ├─ logctx.go
   │  │  │     ├─ README.md
   │  │  │     ├─ slogger.go
   │  │  │     └─ wsl
   │  │  │        └─ wraper.go
   │  │  └─ README.md
   │  └─ webscraper
   │     ├─ cmd
   │     │  └─ scraper
   │     │     └─ main.go
   │     ├─ config
   │     │  ├─ config.go
   │     │  └─ config.yaml
   │     ├─ data
   │     │  └─ ipinfo_lite.mmdb
   │     ├─ go.mod
   │     ├─ go.sum
   │     ├─ info.json
   │     ├─ internal
   │     │  ├─ app
   │     │  │  └─ scraper.go
   │     │  ├─ bootstrap
   │     │  │  └─ runtime.go
   │     │  ├─ cli
   │     │  │  └─ parser.go
   │     │  ├─ domain
   │     │  │  └─ models.go
   │     │  ├─ infrastructure
   │     │  │  ├─ content
   │     │  │  │  ├─ http_strategy.go
   │     │  │  │  └─ trafilatura_strategy.go
   │     │  │  ├─ dns
   │     │  │  │  └─ resolver.go
   │     │  │  ├─ geo
   │     │  │  │  └─ maxmind_provider.go
   │     │  │  ├─ httpclient
   │     │  │  │  └─ client.go
   │     │  │  ├─ logging
   │     │  │  │  └─ logger.go
   │     │  │  ├─ parser
   │     │  │  │  └─ html_cleaner.go
   │     │  │  ├─ quality
   │     │  │  │  └─ simple_evaluator.go
   │     │  │  └─ storage
   │     │  │     └─ file_repository.go
   │     │  └─ usecase
   │     │     └─ scrape_service.go
   │     ├─ logs
   │     ├─ README.md
   │     └─ urls.txt
   └─ python
      ├─ neuro
      │  ├─ BERT
      │  │  ├─ BERT.ipynb
      │  │  ├─ BERT_new.ipynb
      │  │  ├─ requirements.txt
      │  │  └─ train_BERT.py
      │  └─ TODO
      └─ README.md

```