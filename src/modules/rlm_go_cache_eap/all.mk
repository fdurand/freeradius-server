#.PHONY: go_clean_rlm_cache_eap go_build_static_rlm_go_cache_eap go_build_dynamic_rlm_go_cache_eap

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)

MODULE  = rlm_go_cache_eap
BASE_DIR_RLM_GO_CACHE_EAP = ${top_srcdir}/src/modules/$(MODULE)
PACKAGE_GOPATH = $(BASE_DIR_RLM_GO_CACHE_EAP)/.gopath
GOPATH   += $(PACKAGE_GOPATH)
LOCAL_GOPATH = $(subst $(SPACE),:,$(GOPATH))
BASE_RLM_GO_CACHE_EAP     = $(PACKAGE_GOPATH)/src/$(MODULE)

TARGET = $(MODULE).a

# Very dirty hack, to prevent ar from trying to link stuff together. cgo has already taken care of that
TGT_LINKER := "echo"

CGO_CFLAGS1 := $(subst -I,-I${top_srcdir}/,$(CFLAGS))
CGO_CFLAGS2 := $(subst -include ,-include ${top_srcdir}/, $(CGO_CFLAGS1))
# This is probably really hacky, but cgo seems to generate code causing errors
CGO_CFLAGS := $(patsubst -W%,,$(CGO_CFLAGS2))

CGO_LDFLAGS := -Wl,--unresolved-symbols=ignore-all
CGO_LDFLAGS += -L${top_srcdir}/build/lib/local/.libs -lfreeradius-radius -lfreeradius-server
CGO_LDFLAGS += $(LDFLAGS)

create_fake_la_files_rlm_go_cache_eap:
	@echo "/usr/local/lib" > $(top_builddir)/$(BUILD_DIR)/lib/$(MODULE).la
	@echo "$(top_builddir)/$(BUILD_DIR)/lib/local/.libs"  > $(top_builddir)/$(BUILD_DIR)/lib/local/$(MODULE).la

go_build_static_rlm_go_cache_eap: $(BASE_RLM_GO_CACHE_EAP) go_build_dynamic_rlm_go_cache_eap create_fake_la_files_rlm_go_cache_eap
	cd $(BASE_RLM_GO_CACHE_EAP) && \
	GOPATH='$(LOCAL_GOPATH)' CGO_CFLAGS='$(CGO_CFLAGS)' CGO_LDFLAGS='$(CGO_LDFLAGS)' \
	go build -buildmode=c-archive \
	-o $(top_builddir)/$(BUILD_DIR)/lib/local/.libs/$(MODULE).a ./ && \
	cp $(top_builddir)/$(BUILD_DIR)/lib/local/.libs/$(MODULE).a $(top_builddir)/$(BUILD_DIR)/lib/.libs/$(MODULE).a

go_build_dynamic_rlm_go_cache_eap: $(BASE_RLM_GO_CACHE_EAP)
	@echo "CGO_CFLAGS $(CGO_CFLAGS)"
	cd $(BASE_RLM_GO_CACHE_EAP) && \
	GOPATH='$(LOCAL_GOPATH)' CGO_CFLAGS='$(CGO_CFLAGS)' CGO_LDFLAGS='$(CGO_LDFLAGS)' \
	go build -buildmode=c-shared \
	-o $(top_builddir)/$(BUILD_DIR)/lib/local/.libs/$(MODULE).so ./ && \
	cp $(top_builddir)/$(BUILD_DIR)/lib/local/.libs/$(MODULE).so $(top_builddir)/$(BUILD_DIR)/lib/.libs/$(MODULE).so

$(BASE_RLM_GO_CACHE_EAP):
	@mkdir -p $(dir $@)
	@ln -sf $(BASE_DIR_RLM_GO_CACHE_EAP) $@

build/lib/local/rlm_go_cache_eap.la: go_build_static_rlm_go_cache_eap

build/lib/local/rlm_go_cache_eap.so: go_build_dynamic_rlm_go_cache_eap

clean_rlm_go_cache_eap.la: go_clean_rlm_go_cache_eap

go_clean_rlm_go_cache_eap:
	@rm -rf $(GOPATH)
