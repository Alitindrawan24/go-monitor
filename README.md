# Go Monitor

## Overview
This application is designed to monitor the status of websites and send notifications when a website goes down. It reads the monitoring setup from a JSON file, periodically checks the status of the configured websites, and sends notifications via webhooks if any of them are unreachable.

## Installation
1. Clone this repository: `git clone https://github.com/alitindrawan24/go-monitor.git`
2. Navigate to the project directory: `cd go-monitor`
3. Copy the `setup-example.json` as `setup.json` and fill with the desired monitoring configuration.
4. Run the command `go run main.go` or `make run`.
5. The application will start monitoring the specified websites.
6. Notifications will be sent to the configured webhooks if any website goes down.

## Configuration
The monitoring setup is specified in a JSON file named `setup.json` and can be copied from `setup-example.json`, which should be located in the same directory as the executable. The JSON file should have the following structure:
- `interval`: Interval in minutes between each check.
- `targets`: List of URLs to monitor.
- `notification_webhooks`: List of URLs for sending notifications.

Ensure that the JSON file is properly formatted according to the specified structure.

## Dependencies
This application uses Go's standard library for HTTP requests and JSON parsing, with no external dependencies.

## Limitations
While this application provides basic website monitoring functionality, there are some limitations to consider:
- **Goroutine Management**: The application spawns a goroutine for each website check but does not explicitly manage the lifecycle of these goroutines. In scenarios where websites are frequently added or removed, or if the application runs for an extended period, this could lead to potential goroutine leaks and increased memory consumption. Proper goroutine management strategies, such as limiting concurrency or using a worker pool pattern, should be considered to mitigate this risk.
- **Resource Intensive**: Running the application for monitoring a large number of websites continuously may consume significant system resources, especially if many websites are frequently checked.
- **No Historical Data**: The application does not store historical data about website uptime or downtime. It only provides real-time notifications when a website goes down.
- **No Authentication**: The application does not currently support authentication for accessing the monitoring dashboard or API endpoints. Ensure that the application is deployed in a secure environment to prevent unauthorized access.
- **Limited Notification Channels**: The application supports sending notifications only via webhooks. Additional notification channels, such as email or SMS, are not currently implemented.
- **No Customizable Alerts**: The format and content of notification messages are fixed and not customizable. Users cannot configure custom alert thresholds or define specific actions to be taken upon website failure.
- **JSON Configuration Only**: The monitoring setup is configured using a JSON file. There is no support for alternative configuration methods, such as environment variables or command-line arguments.
- **Basic Error Handling**: While the application includes error handling for common scenarios, it may not handle all edge cases or unexpected errors gracefully. Enhancements to error handling could improve the robustness of the application.

## Contributing
Contributions are welcome! If you have any ideas for improvements or new features, feel free to open an issue or submit a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.