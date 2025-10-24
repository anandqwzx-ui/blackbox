const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080/api/v1';

interface RequestOptions extends RequestInit {
  token?: string;
}

async function request<T>(path: string, { token, headers, ...init }: RequestOptions = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
  });

  const isJSON = response.headers.get('content-type')?.includes('application/json');
  const payload = isJSON ? await response.json() : null;

  if (!response.ok) {
    const message = payload?.error ?? response.statusText;
    throw new Error(message);
  }

  return (payload as T) ?? ({} as T);
}

export interface AuthResponse {
  userUuid: string;
  email: string;
  token: string;
}

export async function signup(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/auth/signup', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export interface OrganisationResponse {
  uuid: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export async function listOrganisations(token: string): Promise<OrganisationResponse[]> {
  return request<OrganisationResponse[]>('/organisations', { token });
}

export async function createOrganisation(token: string, name: string): Promise<OrganisationResponse> {
  return request<OrganisationResponse>('/organisations', {
    method: 'POST',
    token,
    body: JSON.stringify({ name }),
  });
}

export interface ProjectResponse {
  uuid: string;
  organisationUuid: string;
  name: string;
  description?: string | null;
  code: string;
  createdAt: string;
  updatedAt: string;
}

export async function listProjects(token: string): Promise<ProjectResponse[]> {
  return request<ProjectResponse[]>('/projects', { token });
}

export async function createProject(
  token: string,
  organisationUuid: string,
  name: string,
  description: string
): Promise<ProjectResponse> {
  return request<ProjectResponse>('/projects', {
    method: 'POST',
    token,
    body: JSON.stringify({ organisationUuid, name, description }),
  });
}

export async function updateProject(
  token: string,
  projectUuid: string,
  name: string,
  description: string
): Promise<ProjectResponse> {
  return request<ProjectResponse>(`/projects/${projectUuid}`, {
    method: 'PUT',
    token,
    body: JSON.stringify({ name, description }),
  });
}

export async function deleteProject(token: string, projectUuid: string): Promise<void> {
  await request(`/projects/${projectUuid}`, {
    method: 'DELETE',
    token,
  });
}

export interface EndpointResponse {
  uuid: string;
  projectUuid: string;
  name?: string | null;
  method: string;
  path: string;
  responseStatus: number;
  responseBody: string;
  responseHeaders: Record<string, string>;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export async function listEndpoints(token: string, projectUuid: string): Promise<EndpointResponse[]> {
  return request<EndpointResponse[]>(`/projects/${projectUuid}/endpoints`, { token });
}

interface EndpointPayload {
  name?: string;
  method: string;
  path: string;
  responseStatus: number;
  responseBody: string;
  responseHeaders: Record<string, string>;
  enabled?: boolean;
}

export async function createEndpoint(
  token: string,
  projectUuid: string,
  payload: EndpointPayload
): Promise<EndpointResponse> {
  return request<EndpointResponse>(`/projects/${projectUuid}/endpoints`, {
    method: 'POST',
    token,
    body: JSON.stringify(payload),
  });
}

export async function updateEndpoint(
  token: string,
  endpointUuid: string,
  payload: EndpointPayload
): Promise<EndpointResponse> {
  return request<EndpointResponse>(`/endpoints/${endpointUuid}`, {
    method: 'PUT',
    token,
    body: JSON.stringify(payload),
  });
}

export async function deleteEndpoint(token: string, endpointUuid: string): Promise<void> {
  await request(`/endpoints/${endpointUuid}`, {
    method: 'DELETE',
    token,
  });
}
