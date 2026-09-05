// Sample Unity PlayerController script
// Attach to a GameObject (e.g., Capsule). Set serverUrl in the inspector to your server (http://localhost:8080)

using System.Collections;
using UnityEngine;
using UnityEngine.Networking;

public class PlayerController : MonoBehaviour
{
    public string serverUrl = "http://localhost:8080";
    public float moveSpeed = 5f;
    private Vector3 targetPos;

    void Start()
    {
        targetPos = transform.position;
    }

    void Update()
    {
        // Keyboard movement (for PC testing)
        float h = Input.GetAxis("Horizontal");
        float v = Input.GetAxis("Vertical");
        Vector3 move = new Vector3(h, 0, v) * moveSpeed * Time.deltaTime;
        if (move.sqrMagnitude > 0.0001f)
        {
            transform.position += move;
            targetPos = transform.position;
            // send a simple move validation to the server (non-blocking)
            StartCoroutine(SendMove(transform.position, move.magnitude / Time.deltaTime));
        }

        // Simple touch input (mobile) - tap to move (placeholder)
        if (Input.touchCount > 0)
        {
            Touch t = Input.GetTouch(0);
            if (t.phase == TouchPhase.Began)
            {
                Ray r = Camera.main.ScreenPointToRay(t.position);
                if (Physics.Raycast(r, out RaycastHit hit))
                {
                    targetPos = hit.point;
                }
            }
            // move toward targetPos
            Vector3 dir = (targetPos - transform.position);
            dir.y = 0;
            if (dir.magnitude > 0.1f)
            {
                Vector3 step = dir.normalized * moveSpeed * Time.deltaTime;
                transform.position += step;
                StartCoroutine(SendMove(transform.position, step.magnitude / Time.deltaTime));
            }
        }
    }

    IEnumerator SendMove(Vector3 pos, float speed)
    {
        var req = new MoveRequest { playerId = "player1", x = pos.x, y = pos.y, z = pos.z, speed = speed, ts = System.DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() };
        string json = JsonUtility.ToJson(req);
        using (UnityWebRequest uwr = UnityWebRequest.Post(serverUrl + "/move", ""))
        {
            byte[] bodyRaw = System.Text.Encoding.UTF8.GetBytes(json);
            uwr.uploadHandler = new UploadHandlerRaw(bodyRaw);
            uwr.downloadHandler = new DownloadHandlerBuffer();
            uwr.SetRequestHeader("Content-Type", "application/json");
            yield return uwr.SendWebRequest();
            if (uwr.result != UnityWebRequest.Result.Success)
            {
                Debug.LogWarning("Move request failed: " + uwr.error);
            }
            else
            {
                Debug.Log("Move response: " + uwr.downloadHandler.text);
            }
        }
    }

    [System.Serializable]
    public class MoveRequest
    {
        public string playerId;
        public float x;
        public float y;
        public float z;
        public float speed;
        public long ts;
    }
}
