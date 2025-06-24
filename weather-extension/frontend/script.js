const backendUrl = "http://localhost:8080/weather";
const tempEl = document.getElementById("temp");
const locationEl = document.getElementById("location");
const weatherIconEl = document.getElementById("weatherIcon");
const settingsPanel = document.getElementById("settingsPanel");

function getWeatherIcon(desc) {
  const lower = desc.toLowerCase();
  if (lower.includes("sun") || lower.includes("clear")) return "☀️";
  if (lower.includes("cloud")) return "☁️";
  if (lower.includes("rain")) return "🌧️";
  if (lower.includes("thunder")) return "⛈️";
  if (lower.includes("snow")) return "❄️";
  if (lower.includes("mist") || lower.includes("fog")) return "🌫️";
  return "🌡️";
}

function updateWidget(data) {
  tempEl.innerText = `${data.temperature}°C`;
  locationEl.innerText = data.location;
  weatherIconEl.innerText = getWeatherIcon(data.description);
}

function fetchWeather(lat, lon) {
  const url = `${backendUrl}?lat=${lat}&lon=${lon}`;
  fetch(url)
    .then(res => res.json())
    .then(data => updateWidget(data))
    .catch(() => {
      tempEl.innerText = "--°C";
      locationEl.innerText = "Unavailable";
    });
}

function fetchByCity() {
  const city = document.getElementById("cityInput").value.trim();
  if (!city) return;
  fetch(`${backendUrl}?city=${encodeURIComponent(city)}`)
    .then(res => res.json())
    .then(data => {
      updateWidget(data);
      settingsPanel.style.display = "none";
    });
}

navigator.geolocation.getCurrentPosition(
  pos => fetchWeather(pos.coords.latitude, pos.coords.longitude),
  err => {
    locationEl.innerText = "Enter City";
    weatherIconEl.innerText = "🌍";
  }
);

// toggle settings panel on click
document.getElementById("menuBtn").onclick = () => {
  settingsPanel.style.display =
    settingsPanel.style.display === "none" ? "flex" : "none";
};

document.getElementById("updateBtn").addEventListener("click", fetchByCity);
