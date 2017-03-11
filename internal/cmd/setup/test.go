package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/therecipe/qt/internal/cmd/deploy"
	"github.com/therecipe/qt/internal/cmd/minimal"
	"github.com/therecipe/qt/internal/cmd/moc"

	"github.com/therecipe/qt/internal/utils"
)

func Test(target string) {
	utils.Log.Infof("running: 'qtsetup test %v'", target)

	if utils.CI() {
		if !(runtime.GOOS == "windows" || target == "windows") { //TODO: split test for windows ?
			utils.Log.Infof("running setup/test %v CI", target)

			path := utils.GoQtPkgPath("internal", "cmd", "moc", "test")

			moc.QmakeMoc(path, target)
			minimal.QmakeMinimal(path, target)

			cmd := exec.Command("go", "test", "-v", "-tags=minimal")
			cmd.Dir = path
			if runtime.GOOS == "windows" {
				for key, value := range map[string]string{
					"PATH":   os.Getenv("PATH"),
					"GOPATH": utils.MustGoPath(),
					"GOROOT": runtime.GOROOT(),

					"TMP":  os.Getenv("TMP"),
					"TEMP": os.Getenv("TEMP"),

					"GOOS":   runtime.GOOS,
					"GOARCH": "386",

					"CGO_ENABLED": "1",
				} {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", key, value))
				}
			}
			utils.RunCmd(cmd, "run qtmoc")
		}
	}

	mode := "test"
	var examples map[string][]string
	if utils.CI() {
		mode = "build"
		examples = map[string][]string{
			"androidextras": []string{"jni", "notification"},

			"canvas3d": []string{"framebuffer", "interaction", "jsonmodels",
				"quickitemtexture", "textureandlight",
				filepath.Join("threejs", "cellphone"),
				filepath.Join("threejs", "oneqt"),
				filepath.Join("threejs", "planets"),
			},

			//"charts": []string{"audio"}, //TODO: ios, ios-simulator

			//"grpc": []string{"hello_world","hello_world2"},

			"gui": []string{"analogclock", "rasterwindow"},

			"qml": []string{"application", "drawer_nav_x", "gallery", "material",
				"prop", "prop2" /*"webview"*/},

			"qt3d": []string{"audio-visualizer-qml"},

			"quick": []string{"bridge", "bridge2", "calc", "dialog", "dynamic",
				"hotreload", "listview", "sailfish", "tableview", "translate", "view"},

			"sailfish": []string{"listview", "listview_variant"},

			"sql": []string{"querymodel"},

			"uitools": []string{"calculator"},

			"widgets": []string{"bridge2" /*"dropsite"*/, "graphicsscene", "line_edits", "pixel_editor",
				/*"renderer"*/ "systray" /*"table"*/, "textedit", filepath.Join("treeview", "treeview_dual"),
				filepath.Join("treeview", "treeview_filelist"), "video_player" /*"webengine"*/},
		}
	} else {
		if strings.HasPrefix(target, "sailfish") {
			examples = map[string][]string{
				"quick": []string{"sailfish"},

				"sailfish": []string{"listview", "listview_variant"},
			}
		} else {
			examples = map[string][]string{
				"qml": []string{"application", "drawer_nav_x", "gallery"},

				"quick": []string{"calc"},

				"widgets": []string{"line_edits", "pixel_editor", "textedit"},
			}
		}
	}

	for cat, list := range examples {
		for _, example := range list {
			if target != "desktop" && example == "textedit" {
				continue
			}
			example := filepath.Join(cat, example)
			utils.Log.Infoln("testing", example)
			deploy.Deploy(&deploy.State{
				BuildMode:   mode,
				BuildTarget: strings.TrimSuffix(target, "-docker"),
				AppPath:     utils.GoQtPkgPath("internal", "examples", example),
				BuildDocker: strings.HasSuffix(target, "-docker"),
			})
		}
	}
}
