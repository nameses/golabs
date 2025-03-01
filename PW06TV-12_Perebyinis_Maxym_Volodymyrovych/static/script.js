document.addEventListener("DOMContentLoaded", function () {
    fetch("/results")
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector("#results-table tbody");
            const totalResultsDiv = document.getElementById("total-results");

            // Populate table with individual results
            data.individualResults.forEach(result => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${result.Name}</td>
                    <td>${result.NP}</td>
                    <td>${result.NPK.toFixed(2)}</td>
                    <td>${result.NPKtg.toFixed(2)}</td>
                    <td>${result.NP2}</td>
                    <td>${result.Ip.toFixed(2)}</td>
                `;
                tableBody.appendChild(row);
            });

            // Populate total results
            totalResultsDiv.innerHTML = `
                <p><strong>Загальний результат ШР1, ШР2, ШР3:</strong></p>
                <p>Кількість: ${data.totalDistributionTiresResult.N}</p>
                <p>n * P: ${data.totalDistributionTiresResult.NP}</p>
                <p>Коефіцієнт використання: ${data.totalDistributionTiresResult.K.toFixed(2)}</p>
                <p>n * P * K: ${data.totalDistributionTiresResult.NPK.toFixed(2)}</p>
                <p>n * P * K * tg: ${data.totalDistributionTiresResult.NPKtg.toFixed(2)}</p>
                <p>n * P^2: ${data.totalDistributionTiresResult.NP2}</p>
                <p>Ефективна кількість ЕП: ${data.totalDistributionTiresResult.NEf.toFixed(2)}</p>
                <p>Розрахунковий коефіцієнт активної потужності: ${data.totalDistributionTiresResult.Ka.toFixed(2)}</p>
                <p>Розрахункове активне навантаження: ${data.totalDistributionTiresResult.Ra.toFixed(2)}</p>
                <p>Розрахункове реактивне навантаження: ${data.totalDistributionTiresResult.Rr.toFixed(2)}</p>
                <p>Повна потужність: ${data.totalDistributionTiresResult.Sr.toFixed(2)}</p>
                <p>Розрахунковий струм: ${data.totalDistributionTiresResult.Ip.toFixed(2)}</p>
            `;
        });
});