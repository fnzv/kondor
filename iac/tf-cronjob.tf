resource "kubernetes_cron_job_v1" "kondor" {
  metadata {
    name = "kondor"
  }
  spec {
    concurrency_policy            = "Replace"
    failed_jobs_history_limit     = 0
    schedule                      = "* * * * *"
    starting_deadline_seconds     = 10
    successful_jobs_history_limit = 2
    job_template {
      metadata {}
      spec {
        backoff_limit              = 2
        ttl_seconds_after_finished = 10
        template {
          metadata {}
          spec {
            container {
              name    = "kondor"
              image   = "fnzv/kondor"
	      #vars from ENV TF_VAR_name_of_var
	      env     {
	       name = "MYSQL_CONN"
  	       value = var.nvr_notifier_mysql
		}
		env {
	       name = "TGBOT_CHATID"
  	       value = var.nvr_notifier_tgbot_chatid
		}
		env {
	       name = "TGBOT_TOKEN"
  	       value = var.nvr_notifier_tgbot_token
		}
		env {
	       name = "FRIGATE_URL"
  	       value = var.nvr_notifier_frigate_url
		}
              command = ["./kondor"]
            }
          }
        }
      }
    }
  }
}
