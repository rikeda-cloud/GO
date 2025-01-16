class Tags {
	constructor() {
		this.tagSelect = document.getElementById("tagSelect");
		this.customTagInput = document.getElementById("customTagInput");
		this.tagSelect.addEventListener("change", this.tagEventHandler.bind(this))
	}

	tagEventHandler() {
		if (tagSelect.value === "custom") {
			customTagInput.style.display = "inline-block";
		} else {
			customTagInput.style.display = "none";
			customTagInput.value = ""; // カスタム入力をクリア
		}
	}

	// タグ情報を取得する関数
	getSelectedTag() {
		if (tagSelect.value === "custom") {
			return customTagInput.value;
		}
		return tagSelect.value;
	}
}
