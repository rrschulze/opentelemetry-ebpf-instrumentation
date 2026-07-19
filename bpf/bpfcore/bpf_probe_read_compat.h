/* SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause) */
/* Copyright The OpenTelemetry Authors */

/*
 * s390x compat shim for bpf_probe_read.
 *
 * On s390x kernels >= 5.5 the BPF verifier no longer permits helper #4
 * (bpf_probe_read) in kprobe programs.  Two typed replacements exist:
 *   - bpf_probe_read_kernel (#113) for kernel virtual addresses
 *   - bpf_probe_read_user   (#112) for user virtual addresses
 *
 * We cannot blindly remap bpf_probe_read to one or the other because the
 * codebase uses it for both kernel structs and user-space buffers.  Instead
 * we leave bpf_probe_read pointing at #113 (kernel) by default for s390x —
 * which covers the large majority of call sites — and individually replace
 * the user-space read sites with explicit bpf_probe_read_user calls.
 *
 * The user-space read sites that have been converted are:
 *   bpf/watcher/watcher.c              (PT_REGS_PARM2 of sys_bind)
 *   bpf/common/sockaddr.h              (get_sockaddr_port_user, parse_sockaddr_info)
 *   bpf/common/iov_iter.h              (ctx->ubuf)
 *   bpf/generictracer/protocol_tcp.h   (u_buf socket payload reads)
 *   bpf/generictracer/protocol_http.h  (args->u_buf)
 *   bpf/generictracer/protocol_http2.h (u_buf)
 *   bpf/generictracer/protocol_handler.c (args->u_buf)
 *   bpf/generictracer/k_tracer_defs.h  (args->u_buf)
 *   bpf/generictracer/java_tls.c       (PT_REGS_PARM* of syscall wrapper)
 *   bpf/generictracer/k_tracer.c       (PT_REGS_PARM1)
 */

#ifndef __BPF_PROBE_READ_COMPAT_H__
#define __BPF_PROBE_READ_COMPAT_H__

#ifdef __TARGET_ARCH_s390
/* Redirect the generic helper to the kernel-memory variant. */
#define bpf_probe_read bpf_probe_read_kernel
#endif

#endif /* __BPF_PROBE_READ_COMPAT_H__ */
