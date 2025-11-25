document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".add-to-cart-form").forEach(form => {
        form.addEventListener("submit", async (e) => {
            e.preventDefault(); // stop redirect

            const formData = new FormData(form);

            const response = await fetch("/cart/add", {
                method: "POST",
                body: formData
            });

            if (response.ok) {
                showPopup("Item added to cart!");
            } else {
                showPopup("Error adding item!");
            }
        });
    });
});

function showPopup(message) {
    const popup = document.createElement("div");
    popup.textContent = message;

    popup.style.position = "fixed";
    popup.style.top = "30px";
    popup.style.right = "30px";
    popup.style.background = "#4CAF50";
    popup.style.color = "white";
    popup.style.padding = "12px 20px";
    popup.style.borderRadius = "8px";
    popup.style.fontSize = "16px";
    popup.style.boxShadow = "0 4px 10px rgba(0,0,0,0.25)";
    popup.style.zIndex = "9999";
    popup.style.opacity = "1";
    popup.style.transition = "opacity 0.5s";

    document.body.appendChild(popup);

    setTimeout(() => {
        popup.style.opacity = "0";
        setTimeout(() => popup.remove(), 500);
    }, 1500);
}
