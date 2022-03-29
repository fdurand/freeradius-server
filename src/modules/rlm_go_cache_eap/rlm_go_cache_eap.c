#include <freeradius-devel/radiusd.h>
#include <freeradius-devel/modules.h>
#include <freeradius-devel/rad_assert.h>

typedef struct rlm_go_cache_eap_t {
       char const      *redis_server;
} rlm_go_cache_eap_t;

static const CONF_PARSER module_config[] = {
         {"redis_server", FR_CONF_OFFSET( PW_TYPE_STRING, rlm_go_cache_eap_t, redis_server), "127.0.0.1" },
         { NULL, -1, 0, NULL, NULL }
 };

extern int go_instantiate(CONF_SECTION *conf, char const *redis_server);
extern int go_authorize(char const *redisServer, REQUEST *request);

static int mod_instantiate(CONF_SECTION *conf, void *instance) {
  rlm_go_cache_eap_t *inst = instance;
  radlog(L_WARN, "Found plugin %s",inst->redis_server);
  return go_instantiate(conf, inst->redis_server);
}

static rlm_rcode_t CC_HINT(nonnull) mod_authorize(void *instance, REQUEST *request) {
  rlm_go_t *inst = instance;
  return go_authorize(inst->redis_server,request);
}

static rlm_rcode_t CC_HINT(nonnull) mod_authenticate(UNUSED void *instance, UNUSED REQUEST *request) {
  return RLM_MODULE_OK;
}

static rlm_rcode_t CC_HINT(nonnull) mod_preacct(UNUSED void *instance, UNUSED REQUEST *request) {
  return RLM_MODULE_OK;
}

static rlm_rcode_t CC_HINT(nonnull) mod_accounting(UNUSED void *instance, UNUSED REQUEST *request) {
    return RLM_MODULE_OK;
}

static rlm_rcode_t CC_HINT(nonnull) mod_checksimul(UNUSED void *instance, REQUEST *request) {
  return RLM_MODULE_OK;
}

static int mod_detach(UNUSED void *instance) {
  return 0;
}

extern module_t rlm_go_cache_eap;
module_t rlm_go = {
  .magic = RLM_MODULE_INIT,
  .name = "go_cache_eap",
  .type = RLM_TYPE_THREAD_SAFE,
  .inst_size = sizeof(rlm_go_cache_eap_t),
  .config = module_config,
  .instantiate = mod_instantiate,
  .detach = mod_detach,
  .methods = {
    [MOD_AUTHENTICATE]      = mod_authenticate,
    [MOD_AUTHORIZE]         = mod_authorize,
#ifdef WITH_ACCOUNTING
    [MOD_PREACCT]           = mod_preacct,
    [MOD_ACCOUNTING]        = mod_accounting,
    [MOD_SESSION]           = mod_checksimul
#endif
  }
};
