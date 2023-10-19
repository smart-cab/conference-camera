import axios from "axios";

const API_URL = "http://localhost:8768/api/v1/";

class AuthService {
  login(email, password) {
    return axios
      .post(API_URL + "login", {
        email,
        password
      })
      .then(response => {
        console.log(response.data)
        if (response.data.data.token) {
          localStorage.setItem("user", JSON.stringify(response.data.data.user));
          localStorage.setItem("token", response.data.data.token)
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem("user");
    localStorage.removeItem("token");
  }

  register(username, email, password) {
    return axios.post(API_URL + "signup", {
      username,
      email,
      password
    });
  }

  getCurrentUser() {
    return JSON.parse(localStorage.getItem('user'));
  }

  getToken() {
    return localStorage.getItem('token');
  }
}

export default new AuthService();
