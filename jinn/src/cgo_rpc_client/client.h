#ifndef CLIENT_H_
#define CLIENT_H_

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

#ifdef CLIENT_EXPORTS
#ifdef _WINDOWS
#define DLL_EXPORT __declspec(dllexport)
#elif defined __linux__
#define DLL_EXPORT __attribute__ ((visibility("default")))
#else
#error "unknown platform"
#endif
#else
#define DLL_EXPORT
#endif // CLIENT_EXPORTS

DLL_EXPORT void* ClientConnect(const char *ip, int port, int max);
DLL_EXPORT void ClientDisconnect(void *handle);
DLL_EXPORT void ClientSetCenterAuth(void *handle, const char *user, const char *passwd, const char *token);
DLL_EXPORT void ClientSetModuleAuth(void *handle, const char *key);
DLL_EXPORT void ClientSendMsg(void *handle, void *message);
DLL_EXPORT void* MessageCreate();
DLL_EXPORT void MessageDestroy(void **handle);
DLL_EXPORT void MessageWriteRequest(void *handle, const char *method,
		const char *params, int params_len, void *data, int data_len);
DLL_EXPORT void MessageReadRespond(void *handle, const char **result, int *result_len, void **data, int *data_len);
DLL_EXPORT void MessageReadError(void *handle, int *code, const char **info);

#ifdef __cplusplus
}
#endif // __cplusplus
#endif // CLIENT_H_
