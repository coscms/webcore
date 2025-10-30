# dbgenerator
if which dbgenerator &> /dev/null; then
    echo "found dbgenerator"
else
    echo "installing dbgenerator"
    go install github.com/webx-top/db/cmd/dbgenerator@latest
fi

DB_PWD="root"
DB_NAME="nging"
dbgenerator -h 127.0.0.1 -d ${DB_NAME} -p ${DB_PWD} -o ./dbschema -match "^(nging_alert_recipient|nging_alert_topic|nging_cloud_backup|nging_cloud_backup_log|nging_cloud_storage|nging_code_invitation|nging_code_verification|nging_config|nging_file|nging_file_embedded|nging_file_moved|nging_file_thumb|nging_kv|nging_login_log|nging_sending_log|nging_task|nging_task_group|nging_task_log|nging_user|nging_user_role|nging_user_role_permission|nging_user_u2f|nging_user_oauth|nging_oauth_app|nging_oauth_agree)$" -backup "./library/setup/install.sql" -charset utf8mb4 #-container "mysql8"
