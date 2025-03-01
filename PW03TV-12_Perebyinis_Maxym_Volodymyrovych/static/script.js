document.getElementById("calc-form").addEventListener("submit", async function (event) {
    event.preventDefault();

    const power = parseFloat(document.getElementById("power").value);
    const deviationCurrent = parseFloat(document.getElementById("deviation_current").value);
    const deviationFinal = parseFloat(document.getElementById("deviation_final").value);
    const price = parseFloat(document.getElementById("price").value);

    const requestData = {
        power: power,
        deviation_current: deviationCurrent,
        deviation_final: deviationFinal,
        price: price
    };

    try {
        const response = await fetch("/calculate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });

        if (!response.ok) {
            throw new Error("Server error: " + response.status);
        }

        const data = await response.json();
        document.getElementById("income_current").textContent = data.income_current.toFixed(2);
        document.getElementById("penalty_current").textContent = data.penalty_current.toFixed(2);
        document.getElementById("income_after").textContent = data.income_after.toFixed(2);
        document.getElementById("penalty_after").textContent = data.penalty_after.toFixed(2);
        document.getElementById("income_final").textContent = data.income_final.toFixed(2);
        
        document.getElementById("results").classList.remove("hidden");
    } catch (error) {
        console.error("Error:", error);
        alert("Failed to calculate. Please try again.");
    }
});
