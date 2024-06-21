const parent = document.querySelector(".entries");

const fetchDataAndRender = async () => {
    parent.innerHTML = '<button class="new">+</button>';
    try {
        const response = await fetch("/getinventory");
        const data = await response.json();

        data.forEach((entry, index) => {
            const container = document.createElement("div");
            container.className = "entrycontainer";
            container.setAttribute("data-id", entry.ID);
            container.innerHTML = `
                <div class="informationDiv">
                    <div class="entrytitle">
                        <h2 class="titleh2">${entry.Name}</h2>
                    </div>
                    <div class="secondaryInfo">
                        <div class="tags">
                            <p class="tag">${entry.Category}</p>
                        </div>
                        <div class="count">
                            <p class="p">${entry.Quantity}</p>
                        </div>
                    </div>
                </div>
                <div class="actionDiv">
                    <button class="edit" onClick="editButtonPressed(this)">🖍️</button>
                    <button class="remove" onClick="remove(this)">🗑️</button>
                </div>
            `;
            container.dataset.index = index;
            parent.appendChild(container);
        });
    } catch (error) {
        console.error("Error fetching data:", error);
    }
};
fetchDataAndRender();