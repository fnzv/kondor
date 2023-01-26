# Kondor

![](imgs/frigate-logo.svg?raw=true)

Kondor is a simple Frigate NVR Telegram notification tool written in Golang

The goal of this script is to notify you on a specific Telegram chat or group as soon as a Frigate event is detected by your CCTV cameras

What is Frigate?

Frigate in short, is a fantastic NVR opensource software that does also object detection out of the box

If you want to see more details visit their website here:
- https://frigate.video

## Quick start

Run the following docker command to bring up a ready to use Mariadb (MySQL)

```
# docker run --name mariadbfrigate -e MARIADB_DATABASE=frigate_db -e MYSQL_ROOT_PASSWORD=mypass -p 3306:3306 -d docker.io/library/mariadb:10.3
```

Populate the following ENV vars in order to be able to run the notify script and get notification on Telegram:

```
# export MYSQL_CONN="root:mypass@tcp(127.0.0.1)/frigate_db?charset=utf8"
# export TGBOT_CHATID="YOUR_TELEGRAM_CHAT_ID"
# export TGBOT_TOKEN="YOUR_TELEGRAM_BOT_TOKEN"
# export FRIGATE_URL="http://my.frigate.nvr.webui.lan"
```

Now run the notify script

```
# go run kondor.go
```

At this moment you should receive notification in the defined Telegram Chat if any event has been triggered or is available in Frigate history

## Getting started

To see how the application works you can just run the commands below on any linux system, the only requirement is Golang and a MySQL database as leverage

Besure to have those requisited below:
- <b>MySQL</b> database is being used to store events plus notification metadata such as sent/not sent
- <b>Golang</b>

1) The first step is to prepare the required ENV variables to be used `MYSQL_CONN,TGBOT_CHATID,TGBOT_TOKEN,FRIGATE_URL`
```
# export MYSQL_CONN="frigate_user:frigate_pass@tcp(nvr-mysql-mariadb.default.svc.cluster.local)/frigate_db?charset=utf8"
# export TGBOT_CHATID="YOUR_TELEGRAM_CHAT_ID"
# export TGBOT_TOKEN="YOUR_TELEGRAM_BOT_TOKEN"
# export FRIGATE_URL="http://my.frigate.nvr.webui.lan"
```

Notes:<br>
You need to have a running MySQL instance that allows the above connection string, the schemas will be created automatically 
<br>
To find out your Telegram ChatID you can write a message to @raw_data_bot and the bot will tell you your ChatID
<br>
To get a Telegram Bot Token you need to contact @BotFather and create a new bot<br>

2) Now that you have your ENV vars set you can run the script as follows:

```
# go run kondor.go
```


The script will populate inside the mysql configured a new schema which will be used to store already sent notifications as a reference

After that it will go through the Frigate NVR APIs "FRIGATE_URL/api/events" by collecting all the stored events and as soon as a new events is scraped it will be notified with the attached snapshot

## Kubernetes

If you are deploying Kondor inside a K8s cluster you can find the Terraform example file under the `iac/` folder <br>

If you want to deploy Kondor via manifest there is also an example for that, the image used was build and pushed on Dockerhub, in case you do not trust it feel free to rebuild it as you can take the Dockerfile in this repo as an example

## Support
Feel free to open issues on this repo, this code is just a quick hack made up to avoid using other notification systems and be reusable in K8s environments
<br>Any support or reply time strongly depends on free time

## Contributing
Feel free to contribute if you think it may be helpful for others

