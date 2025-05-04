function fetchLoginInfo() {
	if (sessionStorage.getItem("username")) {
		return;
	}

	fetch("/api/user")
		.then(res => {
			if (!res.ok) throw new Error("認証エラー");
			return res.json();
		})
		.then(data => {
			sessionStorage.setItem("username", data.user);
		})
		.catch(() => {
			window.location.href = "/login";
		});
}

fetchLoginInfo();
