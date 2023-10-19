import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://localhost:8768/api/v1/';

class UserService {
  ping() {
    return axios.get(API_URL + 'ping');
  }

  getUserBoard() {
    return axios.get(API_URL + 'user/me', { headers: authHeader() });
  }
}

export default new UserService();
