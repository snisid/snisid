import json
from typing import Optional

import httpx
from .exceptions import (
    SNISIDError, AuthenticationError, PermissionError,
    NotFoundError, RateLimitError, ServerError,
)


class HttpClient:
    def __init__(self, api_key: Optional[str] = None, token: Optional[str] = None,
                 base_url: str = "https://api.sniside.ht", cert: Optional[tuple] = None,
                 timeout: int = 30):
        self.base_url = base_url.rstrip("/")
        headers = {"Content-Type": "application/json", "User-Agent": "sniside-sdk/2.1.0"}
        if api_key:
            headers["X-API-Key"] = api_key
        if token:
            headers["Authorization"] = f"Bearer {token}"

        self.client = httpx.Client(
            headers=headers,
            cert=cert,
            timeout=timeout,
            verify=True,
        )

    def _handle_response(self, resp: httpx.Response) -> dict:
        if resp.status_code == 401:
            raise AuthenticationError(resp.text)
        if resp.status_code == 403:
            raise PermissionError(resp.text)
        if resp.status_code == 404:
            raise NotFoundError(resp.text)
        if resp.status_code == 429:
            raise RateLimitError(resp.text)
        if resp.status_code >= 500:
            raise ServerError(resp.text)
        if resp.status_code >= 400:
            raise SNISIDError(resp.status_code, resp.text)
        return resp.json() if resp.text else {}

    def get(self, path: str, params: dict = None) -> dict:
        resp = self.client.get(f"{self.base_url}{path}", params=params)
        return self._handle_response(resp)

    def post(self, path: str, data: dict = None) -> dict:
        resp = self.client.post(f"{self.base_url}{path}", json=data)
        return self._handle_response(resp)

    def put(self, path: str, data: dict = None) -> dict:
        resp = self.client.put(f"{self.base_url}{path}", json=data)
        return self._handle_response(resp)

    def delete(self, path: str) -> dict:
        resp = self.client.delete(f"{self.base_url}{path}")
        return self._handle_response(resp)


class AsyncHttpClient:
    def __init__(self, api_key: Optional[str] = None, token: Optional[str] = None,
                 base_url: str = "https://api.sniside.ht", cert: Optional[tuple] = None,
                 timeout: int = 30):
        self.base_url = base_url.rstrip("/")
        headers = {"Content-Type": "application/json", "User-Agent": "sniside-sdk/2.1.0"}
        if api_key:
            headers["X-API-Key"] = api_key
        if token:
            headers["Authorization"] = f"Bearer {token}"
        self.headers = headers
        self.cert = cert
        self.timeout = timeout
        self.client = None

    async def _get_client(self):
        if self.client is None:
            self.client = httpx.AsyncClient(
                headers=self.headers, cert=self.cert,
                timeout=self.timeout, verify=True,
            )
        return self.client

    async def _handle_response(self, resp: httpx.Response) -> dict:
        if resp.status_code == 401:
            raise AuthenticationError(resp.text)
        if resp.status_code == 403:
            raise PermissionError(resp.text)
        if resp.status_code == 404:
            raise NotFoundError(resp.text)
        if resp.status_code == 429:
            raise RateLimitError(resp.text)
        if resp.status_code >= 500:
            raise ServerError(resp.text)
        if resp.status_code >= 400:
            raise SNISIDError(resp.status_code, resp.text)
        return resp.json() if resp.text else {}

    async def get(self, path: str, params: dict = None) -> dict:
        client = await self._get_client()
        resp = await client.get(f"{self.base_url}{path}", params=params)
        return await self._handle_response(resp)

    async def post(self, path: str, data: dict = None) -> dict:
        client = await self._get_client()
        resp = await client.post(f"{self.base_url}{path}", json=data)
        return await self._handle_response(resp)

    async def put(self, path: str, data: dict = None) -> dict:
        client = await self._get_client()
        resp = await client.put(f"{self.base_url}{path}", json=data)
        return await self._handle_response(resp)

    async def delete(self, path: str) -> dict:
        client = await self._get_client()
        resp = await client.delete(f"{self.base_url}{path}")
        return await self._handle_response(resp)

    async def close(self):
        if self.client:
            await self.client.aclose()
