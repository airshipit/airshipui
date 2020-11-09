/*
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

export class EncryptionConfig {
    name: string;
    encryptionKeyPath: string;
    decryptionKeyPath: string;
    keySecretName: string;
    keySecretNamespace: string;
}

export class Context {
    name: string;
    contextKubeconf: string;
    manifest: string;
    encryptionConfig: string;
    managementConfiguration: string;
}

export class ContextOptions {
    Name = '';
    Manifest = '';
    ManagementConfiguration = '';
    EncryptionConfig = '';
}

// There's no corresponding ManagementConfigOptions in CTL, so the
// Name property has been deliberately capitalized to make it consistent
// with the other ConfigOptions structs and we can retrieve it without
// special handling
export class ManagementConfig {
    Name = '';
    insecure = false;
    systemActionRetries = 0;
    systemRebootDelay = 0;
    type = '';
    useproxy = false;
}

export class Manifest {
    name: string;
    manifest: CtlManifest;
}

export class CtlManifest {
    phaseRepositoryName: string;
    repositories: object;
    targetPath: string;
    subPath: string;
    metadataPath: string;
}

export class Repository {
    url: string;
    auth: RepoAuth;
    checkout: RepoCheckout;
}

export class RepoAuth {
    type: string;
    keyPass: string;
    sshKey: string;
    httpPass: string;
    sshPass: string;
    username: string;
}

export class RepoCheckout {
    commitHash: string;
    branch: string;
    tag: string;
    remoteRef: string;
    force: boolean;
}

// TODO(mfuller): this isn't currently settable from the CLI
// should we allow it in UI?
export class Permissions {
    DirectoryPermission: number;
    FilePermission: number;
}

export class ManifestOptions {
    Name = '';
    RepoName = '';
    URL = '';
    Branch = '';
    CommitHash = '';
    Tag = '';
    RemoteRef = '';
    Force = false;
    IsPhase = false;
    SubPath = '';
    TargetPath = '';
    MetadataPath = '';
}

export class EncryptionConfigOptions {
    Name = '';
    EncryptionKeyPath = '';
    DecryptionKeyPath = '';
    KeySecretName = '';
    KeySecretNamespace = '';
}
