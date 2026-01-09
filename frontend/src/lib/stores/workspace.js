import { writable } from 'svelte/store';
import api from '$lib/api';
import { browser } from '$app/environment';

function createWorkspaceStore() {
  const { subscribe, set, update } = writable({
    list: [],
    currentId: null,
    loading: false
  });

  async function init() {
    if (!browser) return;
    update((s) => ({ ...s, loading: true }));
    try {
      const res = await api.get('/workspaces');
      const list = res.data || [];
      let currentId = null;

      if (list.length > 0) {
        const storedId = Number(localStorage.getItem('workspace_id') || 0);
        const found = list.find((w) => w.id === storedId);
        currentId = found ? found.id : list[0].id;
      }

      if (currentId) {
        localStorage.setItem('workspace_id', String(currentId));
      }

      set({ list, currentId, loading: false });
    } catch (e) {
      set({ list: [], currentId: null, loading: false });
    }
  }

  async function create(name, description = '') {
    const res = await api.post('/workspaces', { name, description });
    const ws = res.data;
    update((s) => {
      const list = [...s.list, ws];
      localStorage.setItem('workspace_id', String(ws.id));
      return { ...s, list, currentId: ws.id };
    });
    return ws;
  }

  async function updateWorkspace(id, name, description = '') {
    const res = await api.put(`/workspaces/${id}`, { name, description });
    const ws = res.data;
    update((s) => {
      const list = s.list.map((w) => (w.id === id ? ws : w));
      return { ...s, list };
    });
    return ws;
  }

  async function getById(id) {
    const res = await api.get(`/workspaces/${id}`);
    return res.data;
  }

  function setCurrent(id) {
    update((s) => {
      const found = s.list.find((w) => w.id === id);
      if (!found) return s;
      if (browser) {
        localStorage.setItem('workspace_id', String(id));
      }
      return { ...s, currentId: id };
    });
  }

  function getCurrentId() {
    let currentId = null;
    if (browser) {
      const stored = localStorage.getItem('workspace_id');
      if (stored) currentId = Number(stored);
    }
    return currentId;
  }

  return {
    subscribe,
    init,
    create,
    update: updateWorkspace,
    getById,
    setCurrent,
    getCurrentId
  };
}

export const workspace = createWorkspaceStore();




