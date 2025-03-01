function calculate(apiUrl) {
    let data = {};
    document.querySelectorAll("input").forEach(input => {
        data[input.id] = parseFloat(input.value) || 0;
    });

    fetch(apiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        let resultText = "";
        if (apiUrl === "/api/losses") {
            resultText = `Матиматичне сподівання збитків: ${data.lossExpectation.toFixed(2)} грн`;
        } else if (apiUrl === "/api/comparing") {
            resultText = `Частота відмов одноколової системи: ${data.failureRate1.toFixed(3)}\n` +
                         `Частота відмов двоколової системи: ${data.failureRate2.toFixed(3)}\n` +
                         `${data.comparisonText}`;
        }
        document.getElementById("results").textContent = resultText;
    })
    .catch(error => {
        console.error("Fetch error:", error);
        document.getElementById("results").textContent = "Помилка запиту: " + error.message;
    });
}
