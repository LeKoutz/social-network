let ws = null;

export function connectWS() {
    ws = new WebSocket(`ws://${window.location.host}/ws`);
    // debug prints for testing
    ws.onopen = () => console.log('WebSocket connection established');
    ws.onclose = (e) => console.log('WebSocket closed', e.code, e.reason);
    ws.onerror = (e) => console.log('WebSocket error', e);
    ws.onmessage = (e) => console.log('Message received', e.data);
}
