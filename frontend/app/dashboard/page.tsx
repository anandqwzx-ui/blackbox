'use client';

import { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';

import {
  OrganisationResponse,
  ProjectResponse,
  EndpointResponse,
  createOrganisation,
  createProject,
  createEndpoint,
  deleteEndpoint,
  deleteProject,
  listEndpoints,
  listOrganisations,
  listProjects,
  updateEndpoint,
  updateProject,
} from '@/lib/api';

interface EndpointFormState {
  uuid?: string;
  name: string;
  method: string;
  path: string;
  responseStatus: number;
  responseBody: string;
  responseHeaders: string;
  enabled: boolean;
}

const defaultEndpointForm: EndpointFormState = {
  name: '',
  method: 'GET',
  path: '/',
  responseStatus: 200,
  responseBody: '{"message":"ok"}',
  responseHeaders: '{"Content-Type":"application/json"}',
  enabled: true,
};

export default function DashboardPage() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [email, setEmail] = useState('');

  const [organisations, setOrganisations] = useState<OrganisationResponse[]>([]);
  const [projects, setProjects] = useState<ProjectResponse[]>([]);
  const [selectedProjectUuid, setSelectedProjectUuid] = useState<string | null>(null);
  const [endpoints, setEndpoints] = useState<EndpointResponse[]>([]);

  const [orgName, setOrgName] = useState('');
  const [projectName, setProjectName] = useState('');
  const [projectDescription, setProjectDescription] = useState('');
  const [projectOrganisation, setProjectOrganisation] = useState('');

  const [editingProjectUuid, setEditingProjectUuid] = useState<string | null>(null);
  const [editingProjectName, setEditingProjectName] = useState('');
  const [editingProjectDescription, setEditingProjectDescription] = useState('');

  const [endpointForm, setEndpointForm] = useState<EndpointFormState>(defaultEndpointForm);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    const savedToken = localStorage.getItem('mockapi_token');
    const savedEmail = localStorage.getItem('mockapi_email');

    if (!savedToken) {
      router.replace('/login');
      return;
    }

    setToken(savedToken);
    setEmail(savedEmail ?? '');
  }, [router]);

  useEffect(() => {
    if (!token) {
      return;
    }

    const bootstrap = async () => {
      try {
        const [orgs, projs] = await Promise.all([
          listOrganisations(token),
          listProjects(token),
        ]);

        setOrganisations(orgs);
        setProjects(projs);
        if (projs.length > 0) {
          setSelectedProjectUuid((current) => current ?? projs[0].uuid);
        }
      } catch (error) {
        setErrorMessage(error instanceof Error ? error.message : 'Unable to load workspace');
      }
    };

    bootstrap();
  }, [token]);

  useEffect(() => {
    if (!token || !selectedProjectUuid) {
      setEndpoints([]);
      return;
    }

    listEndpoints(token, selectedProjectUuid)
      .then(setEndpoints)
      .catch((error) => setErrorMessage(error instanceof Error ? error.message : 'Unable to load endpoints'));
  }, [token, selectedProjectUuid]);

  const selectedProject = useMemo(
    () => projects.find((project) => project.uuid === selectedProjectUuid) ?? null,
    [projects, selectedProjectUuid]
  );

  const resetEndpointForm = () => {
    setEndpointForm(defaultEndpointForm);
  };

  const handleCreateOrganisation = async () => {
    if (!token || !orgName.trim()) {
      if (!orgName.trim()) {
        setErrorMessage('Organisation name is required');
        setStatusMessage(null);
      }
      return;
    }

    try {
      const organisation = await createOrganisation(token, orgName.trim());
      setOrganisations((current) => [organisation, ...current]);
      setOrgName('');
      setErrorMessage(null);
      setStatusMessage('Organisation created');
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to create organisation');
    }
  };

  const handleCreateProject = async () => {
    if (!token || !projectOrganisation) {
      setErrorMessage('Choose an organisation for the project');
      setStatusMessage(null);
      return;
    }

    if (!projectName.trim()) {
      setErrorMessage('Project name is required');
      setStatusMessage(null);
      return;
    }

    try {
      const project = await createProject(token, projectOrganisation, projectName.trim(), projectDescription.trim());
      setProjects((current) => [project, ...current]);
      setProjectName('');
      setProjectDescription('');
      setErrorMessage(null);
      setStatusMessage('Project created');
      if (!selectedProjectUuid) {
        setSelectedProjectUuid(project.uuid);
      }
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to create project');
    }
  };

  const beginProjectEdit = (project: ProjectResponse) => {
    setEditingProjectUuid(project.uuid);
    setEditingProjectName(project.name);
    setEditingProjectDescription(project.description ?? '');
  };

  const handleSaveProject = async () => {
    if (!token || !editingProjectUuid) {
      return;
    }

    const trimmedName = editingProjectName.trim();
    if (!trimmedName) {
      setErrorMessage('Project name is required');
      return;
    }

    try {
      const updated = await updateProject(token, editingProjectUuid, trimmedName, editingProjectDescription.trim());
      setProjects((current) => current.map((item) => (item.uuid === updated.uuid ? updated : item)));
      setErrorMessage(null);
      setStatusMessage('Project updated');
      setEditingProjectUuid(null);
      setEditingProjectName('');
      setEditingProjectDescription('');
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to update project');
    }
  };

  const handleDeleteProject = async (projectUuid: string) => {
    if (!token) {
      return;
    }

    try {
      await deleteProject(token, projectUuid);
      setProjects((current) => current.filter((project) => project.uuid !== projectUuid));
      if (selectedProjectUuid === projectUuid) {
        setSelectedProjectUuid(null);
        setEndpoints([]);
      }
      setErrorMessage(null);
      setStatusMessage('Project deleted');
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to delete project');
    }
  };

  const handleSubmitEndpoint = async () => {
    if (!token || !selectedProjectUuid) {
      return;
    }

    let headersObject: Record<string, string> = {};
    try {
      headersObject = JSON.parse(endpointForm.responseHeaders || '{}');
    } catch (error) {
      setErrorMessage('Response headers must be valid JSON');
      return;
    }

    try {
      if (endpointForm.uuid) {
        const updated = await updateEndpoint(token, endpointForm.uuid, {
          name: endpointForm.name,
          method: endpointForm.method,
          path: endpointForm.path,
          responseStatus: endpointForm.responseStatus,
          responseBody: endpointForm.responseBody,
          responseHeaders: headersObject,
          enabled: endpointForm.enabled,
        });
        setEndpoints((current) => current.map((item) => (item.uuid === updated.uuid ? updated : item)));
        setStatusMessage('Endpoint updated');
      } else {
        const created = await createEndpoint(token, selectedProjectUuid, {
          name: endpointForm.name,
          method: endpointForm.method,
          path: endpointForm.path,
          responseStatus: endpointForm.responseStatus,
          responseBody: endpointForm.responseBody,
          responseHeaders: headersObject,
          enabled: endpointForm.enabled,
        });
        setEndpoints((current) => [created, ...current]);
        setStatusMessage('Endpoint created');
      }
      setErrorMessage(null);
      resetEndpointForm();
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save endpoint');
    }
  };

  const handleEditEndpoint = (endpoint: EndpointResponse) => {
    setEndpointForm({
      uuid: endpoint.uuid,
      name: endpoint.name ?? '',
      method: endpoint.method,
      path: endpoint.path,
      responseStatus: endpoint.responseStatus,
      responseBody: endpoint.responseBody,
      responseHeaders: JSON.stringify(endpoint.responseHeaders, null, 2),
      enabled: endpoint.enabled,
    });
  };

  const handleDeleteEndpoint = async (endpointUuid: string) => {
    if (!token) {
      return;
    }

    try {
      await deleteEndpoint(token, endpointUuid);
      setEndpoints((current) => current.filter((endpoint) => endpoint.uuid !== endpointUuid));
      setErrorMessage(null);
      setStatusMessage('Endpoint deleted');
    } catch (error) {
      setStatusMessage(null);
      setErrorMessage(error instanceof Error ? error.message : 'Unable to delete endpoint');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('mockapi_token');
    localStorage.removeItem('mockapi_email');
    router.replace('/login');
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2.5rem' }}>
      <section>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', gap: '1rem' }}>
          <div>
            <h1>Dashboard</h1>
            <p className="lead">Manage organisations, projects, and mock endpoints from a single, focused workspace.</p>
          </div>
          <div className="card" style={{ padding: '1rem 1.25rem', maxWidth: '260px' }}>
            <span className="muted">Signed in as</span>
            <strong>{email}</strong>
            <button onClick={handleLogout} style={{ marginTop: '0.75rem' }}>
              Log out
            </button>
          </div>
        </div>

        {(statusMessage || errorMessage) && (
          <div className="card" style={{ marginTop: '1.5rem', background: errorMessage ? 'rgba(255, 80, 80, 0.12)' : 'rgba(120, 255, 120, 0.12)' }}>
            <strong>{errorMessage ?? statusMessage}</strong>
            <button
              onClick={() => {
                setStatusMessage(null);
                setErrorMessage(null);
              }}
              style={{ alignSelf: 'flex-start', marginTop: '0.5rem', background: 'transparent', color: 'inherit' }}
            >
              Dismiss
            </button>
          </div>
        )}
      </section>

      <section>
        <h2>Organisations</h2>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            handleCreateOrganisation();
          }}
          className="card"
        >
          <label htmlFor="orgName">Organisation name</label>
          <input
            id="orgName"
            placeholder="Studio Inc."
            value={orgName}
            onChange={(event) => setOrgName(event.target.value)}
          />
          <button type="submit">Create organisation</button>
        </form>

        <div className="grid" style={{ marginTop: '1.5rem' }}>
          {organisations.map((organisation) => (
            <div key={organisation.uuid} className="card">
              <strong>{organisation.name}</strong>
              <span className="muted">UUID</span>
              <code>{organisation.uuid}</code>
            </div>
          ))}
          {organisations.length === 0 && <p className="muted">No organisations yet.</p>}
        </div>
      </section>

      <section>
        <h2>Projects</h2>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            handleCreateProject();
          }}
          className="card"
        >
          <label htmlFor="organisation">
            Organisation
            <select
              id="organisation"
              value={projectOrganisation}
              onChange={(event) => setProjectOrganisation(event.target.value)}
            >
              <option value="">Select organisation</option>
              {organisations.map((organisation) => (
                <option key={organisation.uuid} value={organisation.uuid}>
                  {organisation.name}
                </option>
              ))}
            </select>
          </label>

          <label htmlFor="projectName">Project name</label>
          <input
            id="projectName"
            placeholder="Mobile API sandbox"
            value={projectName}
            onChange={(event) => setProjectName(event.target.value)}
          />

          <label htmlFor="projectDescription">Description</label>
          <textarea
            id="projectDescription"
            placeholder="Short description"
            value={projectDescription}
            onChange={(event) => setProjectDescription(event.target.value)}
            rows={3}
          />

          <button type="submit">Create project</button>
        </form>

        <div className="grid" style={{ marginTop: '1.5rem' }}>
          {projects.map((project) => {
            const isActive = selectedProjectUuid === project.uuid;
            const isEditing = editingProjectUuid === project.uuid;
            return (
              <div key={project.uuid} className="card" style={{ gap: '0.75rem' }}>
                {isEditing ? (
                  <>
                    <label>
                      Name
                      <input
                        value={editingProjectName}
                        onChange={(event) => setEditingProjectName(event.target.value)}
                        style={{ fontWeight: 600, fontSize: '1.1rem' }}
                      />
                    </label>
                    <label>
                      Description
                      <textarea
                        value={editingProjectDescription}
                        onChange={(event) => setEditingProjectDescription(event.target.value)}
                        rows={3}
                      />
                    </label>
                    <div style={{ display: 'flex', gap: '0.5rem' }}>
                      <button type="button" onClick={handleSaveProject}>
                        Save
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setEditingProjectUuid(null);
                          setEditingProjectName('');
                          setEditingProjectDescription('');
                        }}
                        style={{ background: 'transparent', border: '1px solid #ffffff', color: '#ffffff' }}
                      >
                        Cancel
                      </button>
                    </div>
                  </>
                ) : (
                  <>
                    <strong style={{ fontSize: '1.2rem' }}>{project.name}</strong>
                    <p className="muted">{project.description || 'No description yet.'}</p>
                  </>
                )}
                <span className="muted">Code: {project.code}</span>
                <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    onClick={() => setSelectedProjectUuid(project.uuid)}
                    style={{
                      background: isActive ? '#ffffff' : 'transparent',
                      color: isActive ? '#050505' : '#ffffff',
                      border: '1px solid #ffffff',
                    }}
                  >
                    View endpoints
                  </button>
                  <button type="button" onClick={() => beginProjectEdit(project)}>
                    Edit
                  </button>
                  <button
                    type="button"
                    onClick={() => handleDeleteProject(project.uuid)}
                    style={{ background: 'transparent', border: '1px solid #ff5f5f', color: '#ff5f5f' }}
                  >
                    Delete
                  </button>
                </div>
              </div>
            );
          })}
          {projects.length === 0 && <p className="muted">No projects available.</p>}
        </div>
      </section>

      <section>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', gap: '1rem' }}>
          <h2>Endpoints</h2>
          {selectedProject && <span className="muted">Project code: {selectedProject.code}</span>}
        </div>

        {selectedProject ? (
          <>
            <form
              onSubmit={(event) => {
                event.preventDefault();
                handleSubmitEndpoint();
              }}
              className="card"
            >
              <label htmlFor="endpointName">Name</label>
              <input
                id="endpointName"
                placeholder="List users"
                value={endpointForm.name}
                onChange={(event) => setEndpointForm((form) => ({ ...form, name: event.target.value }))}
              />

              <div style={{ display: 'grid', gap: '0.75rem', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))' }}>
                <label>
                  Method
                  <select
                    value={endpointForm.method}
                    onChange={(event) => setEndpointForm((form) => ({ ...form, method: event.target.value }))}
                  >
                    {['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'].map((method) => (
                      <option key={method} value={method}>
                        {method}
                      </option>
                    ))}
                  </select>
                </label>

                <label>
                  Path
                  <input
                    placeholder="/users"
                    value={endpointForm.path}
                    onChange={(event) => setEndpointForm((form) => ({ ...form, path: event.target.value }))}
                  />
                </label>

                <label>
                  Status
                  <input
                    type="number"
                    min={100}
                    max={599}
                    value={endpointForm.responseStatus}
                    onChange={(event) =>
                      setEndpointForm((form) => ({ ...form, responseStatus: Number(event.target.value) }))
                    }
                  />
                </label>

                <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <input
                    type="checkbox"
                    checked={endpointForm.enabled}
                    onChange={(event) => setEndpointForm((form) => ({ ...form, enabled: event.target.checked }))}
                    style={{ width: '1rem', height: '1rem' }}
                  />
                  Enabled
                </label>
              </div>

              <label htmlFor="endpointBody">Response body</label>
              <textarea
                id="endpointBody"
                rows={4}
                value={endpointForm.responseBody}
                onChange={(event) => setEndpointForm((form) => ({ ...form, responseBody: event.target.value }))}
              />

              <label htmlFor="endpointHeaders">Response headers (JSON)</label>
              <textarea
                id="endpointHeaders"
                rows={3}
                value={endpointForm.responseHeaders}
                onChange={(event) => setEndpointForm((form) => ({ ...form, responseHeaders: event.target.value }))}
              />

              <div style={{ display: 'flex', gap: '0.75rem' }}>
                <button type="submit">{endpointForm.uuid ? 'Update endpoint' : 'Create endpoint'}</button>
                {endpointForm.uuid && (
                  <button
                    type="button"
                    onClick={resetEndpointForm}
                    style={{ background: 'transparent', border: '1px solid #ffffff', color: '#ffffff' }}
                  >
                    Cancel
                  </button>
                )}
              </div>
            </form>

            <div className="grid" style={{ marginTop: '1.5rem' }}>
              {endpoints.map((endpoint) => (
                <div key={endpoint.uuid} className="card" style={{ gap: '0.5rem' }}>
                  <strong>
                    {endpoint.method} {endpoint.path}
                  </strong>
                  <span className="muted">Status {endpoint.responseStatus}</span>
                  <code style={{ whiteSpace: 'pre-wrap' }}>{endpoint.responseBody}</code>
                  <span className="muted">Headers</span>
                  <code style={{ whiteSpace: 'pre-wrap' }}>{JSON.stringify(endpoint.responseHeaders, null, 2)}</code>
                  <span className="muted">Enabled: {endpoint.enabled ? 'Yes' : 'No'}</span>
                  <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                    <button type="button" onClick={() => handleEditEndpoint(endpoint)}>
                      Edit
                    </button>
                    <button
                      type="button"
                      onClick={() => handleDeleteEndpoint(endpoint.uuid)}
                      style={{ background: 'transparent', border: '1px solid #ff5f5f', color: '#ff5f5f' }}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
              {endpoints.length === 0 && <p className="muted">No endpoints configured.</p>}
            </div>
          </>
        ) : (
          <p className="muted">Select a project to manage its endpoints.</p>
        )}
      </section>
    </div>
  );
}
