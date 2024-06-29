export const callApi = async (url, options, method) => {
  // Helper to get the token from local storage
  const getToken = () => localStorage.getItem("access_token");
  const getSession = () => localStorage.getItem("session");

  const fetchData = async () => {
    const response = await fetch(url, {
      ...options,
      method: method,
      headers: {
        ...options.headers,
        Authorization: getToken(),
      },
    });
    return response;
  };

  let response = await fetchData();

  if (response.status === 400) {
    const refreshResponse = await fetch(process.env.NEXT_PUBLIC_REFRESH_TOKEN, {
      method: "POST",
      headers: {
        ...options.headers,
        Authorization: localStorage.getItem("refresh_token"),
      },
    });

    if (refreshResponse.ok) {
      const data = await refreshResponse.json();
      localStorage.setItem("access_token", data.access_token);
      response = await fetchData(); // Retry the request with the new token
    } else {
      const errorData = await refreshResponse.json();
      localStorage.clear();
      window.location.replace("/");
      throw new Error(`Refresh failed: ${errorData.message}`);
    }
  }

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(`${errorData.message}`);
  }

  return await response.json();
};

export const callLogin = async (eamilId, password) => {
  const body = JSON.stringify({
    email: eamilId,
    password: password,
  });
  const payload = {
    method: "POST",
    body: body,
  };
  const data = await fetch(process.env.NEXT_PUBLIC_SIGN_IN, payload);
  const resp = await data.json();
  console.log(resp, resp.access_token, resp.refresh_token, resp.session);
  if (resp.status == 200) {
    localStorage.setItem("email", eamilId);
    localStorage.setItem("access_token", resp.access_token);
    localStorage.setItem("refresh_token", resp.refresh_token);
    localStorage.setItem("session", resp.session);
    localStorage.setItem("type", resp.type);
  }
  return resp;
};
