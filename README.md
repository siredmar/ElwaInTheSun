# Elwa In The Sun

This repository contains a **Go-based backend** that interacts with a **Sonnen Battery** and a **MyPV water heating system**. It provides an HTTP API for control and monitoring, and a simple web frontend for configuring system settings.

## Features

- Communicates with **Sonnen Batteries** to retrieve and manage power data.
- Controls **MyPV water heating systems** based on available PV energy.
- Provides a **REST API** for configuration and monitoring.
- Includes a **lightweight web UI** for setting adjustments.
- Runs as a **Docker container** for easy deployment.

## Configuration

The backend reads its configuration from a JSON file. Below is an example:

```json
{
  "sonnen_token": "mySonnenToken",
  "sonnen_host": "http://192.168.1.123",
  "mypv_token": "myPVToken",
  "mypv_serial": "1111111111111111",
  "interval": "55s",
  "reserved": 100,
  "max_temp": 50
}
```

### Configuration Parameters

- **sonnen_token** – API token for Sonnen Battery integration.
- **sonnen_host** – IP or hostname of the Sonnen Battery API.
- **mypv_token** – API token for MyPV device.
- **mypv_serial** – Serial number of the MyPV device.
- **interval** – Control loop execution interval (e.g., `55s`, `5m`).
- **reserved** – Reserved power (W) to be left in the battery.
- **max_temp** – Maximum water temperature allowed.

## Running with Docker Compose

The system consists of two services:
1. **Controller** - Runs the backend logic and manages the battery & heating system.
2. **Frontend** - Serves the web UI and API for configuration.

### Example `docker-compose.yml`

```yaml
version: "3.3"
services:
  frontend:
    restart: always
    image: siredmar/elwainthesun:v0.3.1
    command: /controller serve -c /config.json
    ports:
      - "8080:8080"
    depends_on:
      - controller
    volumes:
      - /root/elwa/config.json:/config.json
    environment:
      LOG_LEVEL: info
  controller:
    restart: always
    image: siredmar/elwainthesun:v0.3.1
    command: /controller run -c /config.json
    volumes:
      - /root/elwa/config.json:/config.json
    environment:
      LOG_LEVEL: info
```

### Running the Containers

Start the services with:

```sh
docker-compose up -d
```

Check logs:

```sh
docker-compose logs -f
```

### Web Access

Once running, the frontend is accessible at:

```
http://localhost:8080
```

Replace `localhost` with the server’s IP if running remotely.

## API Endpoints

- `GET /settings` – Retrieve current configuration.
- `POST /settings` – Update configuration (expects JSON body).
- `GET /status` – Get system status and power data.

## Deployment Notes

- Ensure that `/root/elwa/config.json` exists before starting the containers.
- The **frontend** container serves the UI and settings API.
- The **controller** container interacts with Sonnen and MyPV devices.

## Troubleshooting

- Check if the containers are running:

  ```sh
  docker ps
  ```

- Verify logs for errors:

  ```sh
  docker-compose logs controller
  ```

- Confirm that `config.json` is correctly mounted:

  ```sh
  docker exec -it <container_id> cat /config.json
  ```

## License

This project is licensed under the MIT License.

