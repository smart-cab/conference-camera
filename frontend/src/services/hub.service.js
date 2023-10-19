import axios from 'axios';

const API_URL = 'http://localhost:8768/api/v1/';

class HubService {
  add(token) {
    return axios
      .get(API_URL + "token", {token})
  }
}

export default new HubService();
