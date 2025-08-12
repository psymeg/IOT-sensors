#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <ArduinoJson.h>
#include "DHT.h"

// ====== CONFIG ======
const char* ssid = "YOUR_WIFI_SSID";
const char* password = "YOUR_WIFI_PASSWORD";
const char* serverUrl = "http://192.168.1.50:8080/sensor"; // Change to your Go server IP

#define DHTPIN D4      // Pin where your DHT22 data line is connected
#define DHTTYPE DHT22  // Sensor type
// ====================

DHT dht(DHTPIN, DHTTYPE);

void setup() {
  Serial.begin(115200);
  WiFi.begin(ssid, password);

  Serial.print("Connecting to WiFi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected!");
  
  dht.begin();
}

void loop() {
  if (WiFi.status() == WL_CONNECTED) {
    float temp = dht.readTemperature();

    if (isnan(temp)) {
      Serial.println("Failed to read from DHT sensor!");
      delay(2000);
      return;
    }

    sendTemperature(temp);
  } else {
    Serial.println("WiFi disconnected!");
    WiFi.begin(ssid, password);
  }

  delay(10000); // Send data every 10 seconds
}

void sendTemperature(float temperature) {
  WiFiClient client;
  HTTPClient http;

  http.begin(client, serverUrl);
  http.addHeader("Content-Type", "application/json");

  // Create JSON payload
  StaticJsonDocument<100> doc;
  doc["temperature"] = temperature;

  String requestBody;
  serializeJson(doc, requestBody);

  int httpCode = http.POST(requestBody);

  if (httpCode > 0) {
    Serial.printf("POST... code: %d\n", httpCode);
    if (httpCode == HTTP_CODE_OK) {
      Serial.println("Data sent successfully");
    }
  } else {
    Serial.printf("POST failed, error: %s\n", http.errorToString(httpCode).c_str());
  }

  http.end();
}

