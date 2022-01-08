# kubemc

Small tool to manage multiple kubeconfig.

* monitor file change event in `~/.kube/kubemc/`
* merge .yaml/.yml config files to `~/.kube/config`

## Usage

1. Build binary
```
git clone https://github.com/puglao/kubemc.git
cd kubemc
go build .
```

2. Create `~.kube/kubemc` directory
```
mkdir -p ~/.kube/kubemc
```

3. Put kubemc binary to `/usr/local/bin`
```
sudo cp kubemc /usr/local/bin
```

4. Generate .plist file <MacOS only>
```
sed "s|{{ HOME }}|$HOME|g" site.cloudemo.kubemc.plist-example > site.cloudemo.kubemc.plist
```

5. copy load plist <MacOS only>
```
cp site.cloudemo.kubemc.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/site.cloudemo.kubemc.plist
```

5. Copy multi-kubeconfig to `~/.kube/kubemc`
```
cp <kubeconfig files> ~/.kube/kubemc/
```

6. Create/modfiy/rename/delete files in `~/.kube/kubemc` will trigger config merge
```
touch ~/.kube/kubemc/test.yaml
```


## Environment Variable
|Environment Variable|Descrption|default value|
|-|-|-|
|`KUBECONFIG`|location of kubeconfig|`~/.kube/config`|
|`KUBEMC_DIR`|kubemc directory|`~/.kube/kubemc`|
|`KUBEMC_RATELIMIT`|One Merge cannot trigger within given period|`2` second|