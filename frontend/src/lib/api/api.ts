import ky from 'ky'


export const getOrCreateDeviceId = () => {
  let deviceId = localStorage.getItem('device_id');

  if (!deviceId) {
    deviceId = crypto.randomUUID(); 
    
    localStorage.setItem('device_id', deviceId);
  }

  return deviceId;
};






let isRefreshing = false
let refreshPromise: Promise<void> | null = null
let refreshSubscribers: Array<(success: boolean) => void> = []

function subscribe(cb: (success: boolean) => void) {
  refreshSubscribers.push(cb)
}

function notifySubscribers(success: boolean) {
  refreshSubscribers.splice(0).forEach(cb => cb(success))
}

async function refreshToken(): Promise<void> {
  if (!refreshPromise) {
    refreshPromise = ky.post('auth/token/refresh', {
	prefixUrl: 'http://localhost:5000/api/v1',
      headers: {
        'X-Device-ID': getOrCreateDeviceId() || '' 
      },
      credentials: 'include',
    })
    .then(() => {
      notifySubscribers(true)
    })
    .catch(err => {
      notifySubscribers(false)
      throw err
    })
    .finally(() => {
      refreshPromise = null
    })
  }
  return refreshPromise
}

export const api = ky.create({
  prefixUrl: 'http://localhost:5000/api/v1', 
  credentials: 'include',
  hooks: {
    beforeRequest: [
      (request) => {
        const deviceId = getOrCreateDeviceId()
        if (deviceId) {
          request.headers.set('X-Device-ID', deviceId)
        }
      }
    ],
    afterResponse: [
      async (request, options, response) => {
        if (
          response.status !== 401 || 
          request.url.includes('/auth/token/refresh') ||
          (options as any)._retried
        ) {
          return response
        }

        if (!isRefreshing) {
		  isRefreshing = true
          try {
            await refreshToken()
          } catch {
            throw response
          }finally {
			isRefreshing = false
		  }
        } else {
          const success = await new Promise<boolean>(resolve => subscribe(resolve))
          if (!success) {
            throw response
          }
        }

        return api(request, { 
          ...options, 
          _retried: true 
        } as any)
      },
    ],
  },
})
