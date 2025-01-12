// カスタムタグ入力の切り替え
const tagSelect = document.getElementById("tagSelect");
const customTagInput = document.getElementById("customTagInput");

tagSelect.addEventListener("change", () => {
	if (tagSelect.value === "custom") {
		customTagInput.style.display = "inline-block";
	} else {
		customTagInput.style.display = "none";
		customTagInput.value = ""; // カスタム入力をクリア
	}
});

// タグ情報を取得する関数
function getSelectedTag() {
	if (tagSelect.value === "custom") {
		return customTagInput.value;
	}
	return tagSelect.value;
}
