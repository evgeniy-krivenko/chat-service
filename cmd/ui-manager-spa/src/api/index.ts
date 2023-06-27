import { v4 as uuidv4 } from 'uuid';
import axios, { AxiosResponse } from 'axios';
import { redirect } from 'react-router-dom';
import { BASE_URL } from '../const/urls';

const $api = axios.create();

$api.interceptors.request.use(async (req) => {
  const token = localStorage.getItem('token');
  if (token) {
    req.headers.Authorization = `Bearer ${token}`;
  }

  req.headers['X-Request-ID'] = uuidv4();

  return req;
});

$api.interceptors.response.use((response: AxiosResponse) => {
  if (response.data?.error) {
    if (response.data?.error?.code === 401) {
      redirect('/login');
    }
    throw Error(response.data?.error?.message);
  }

  return response;
});

$api.defaults.baseURL = BASE_URL;
$api.defaults.timeout = 3000;

export default $api;
