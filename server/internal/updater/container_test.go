package updater

import "testing"

// 只测试环境变量的解析，不调用 InContainer，避免测试运行器自身的 Docker 环境影响结果。
func TestInContainerEnvFlag(t *testing.T) {
	cases := map[string]bool{
		"1":     true,
		"true":  true,
		"yes":   true,
		"0":     false,
		"false": false,
		"FALSE": false,
		"":      false, // 空 = 交给别的判据
	}
	for v, want := range cases {
		t.Setenv("OVH_IN_CONTAINER", v)
		got := containerEnvEnabled()
		// 空串和显式否定时会落到文件/cgroup 判据上，
		// 本机跑测试时那些都不该命中
		if want && !got {
			t.Errorf("OVH_IN_CONTAINER=%q 应当判为容器内", v)
		}
		if !want && got {
			t.Errorf("OVH_IN_CONTAINER=%q 不该判为容器内（会无故停掉裸机的自更新）", v)
		}
	}
}

// 更新指引必须把"为什么"说清楚，否则用户只会觉得功能坏了。
func TestContainerUpdateHintIsActionable(t *testing.T) {
	for _, must := range []string{"docker compose up -d --build", "docker compose pull", "卷", "可写层"} {
		if !contains(ContainerUpdateHint, must) {
			t.Errorf("更新指引里缺少 %q", must)
		}
	}
}

func contains(h, n string) bool {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return true
		}
	}
	return false
}
