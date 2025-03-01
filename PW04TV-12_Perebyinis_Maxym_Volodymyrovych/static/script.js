function calculate(apiUrl) {
    let data = {};
    document.querySelectorAll("input").forEach(input => {
        data[input.id] = parseFloat(input.value);
    });

    fetch(apiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        let resultText = "";
        if (apiUrl === "/api/calc1") {
            resultText = `Струм КЗ нормальний режим: ${data.im.toFixed(1)}\n` +
                         `Струм КЗ аварійний режим: ${data.imPa.toFixed(1)}\n` +
                         `Економічний переріз: ${data.sEk.toFixed(1)}\n` +
                         `Мінімальний переріз: ${data.sVsS.toFixed(1)}`;
        } else if (apiUrl === "/api/calc2") {
            resultText = `Опора елементу xc: ${data.xc.toFixed(2)}\n` +
                         `Опора елементу xt: ${data.xt.toFixed(2)}\n` +
                         `Початкове діюче значення струму: ${data.ip0.toFixed(1)}`;
        } else if (apiUrl === "/api/calc3") {
            resultText = `Трифазний струм КЗ: ${data.ish3.toFixed(1)}\n` +
                         `Двофазний струм КЗ: ${data.ish2.toFixed(1)}`;
        }
        document.getElementById("results").textContent = resultText;
    });
}