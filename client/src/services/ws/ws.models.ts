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

export interface WsReceiver {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    type: string;
    component: string;

    // This is the method which will need to be implemented in the component to handle the messages
    receiver(message: WsMessage): Promise<void>;
}

// WebsocketMessage is the structure for the json that is used to talk to the backend
export class WsMessage {
  sessionID: string;
  type: string;
  component: string;
  subComponent: string;
  timestamp: number;
  dashboards: Dashboard[];
  error: string;
  html: string;
  name: string;
  details: string;
  id: string;
  isAuthenticated: boolean;
  message: string;
  token: string;
  refreshToken: string;
  data: JSON;
  yaml: string;
  actionType: string;
  targets: string[];
  authentication: Authentication;

  // this constructor looks like this in case anyone decides they want just a raw message with no data predefined
  // or an easy way to specify the defaults
  constructor(type?: string | null, component?: string | null, subComponent?: string | null) {
    this.type = type;
    this.component = component;
    this.subComponent = subComponent;
  }
}

export class WsConstants {
  // CTL constants
  public static readonly BAREMETAL = 'baremetal';
  public static readonly BUNDLE = 'bundle';
  public static readonly CTL = 'ctl';
  public static readonly CLUSTER = 'cluster';
  public static readonly CONFIG = 'config';
  public static readonly DECRYPT = 'decrypt';
  public static readonly DETAILS = 'details';
  public static readonly DOCUMENT = 'document';
  public static readonly DOCUMENTS = 'documents';
  public static readonly ENCRYPT = 'encrypt';
  public static readonly ERROR = 'error';
  public static readonly EXECUTOR = 'executor';
  public static readonly GET_DEFAULTS = 'getDefaults';
  public static readonly GENERATE = 'generate';
  public static readonly IMAGE = 'image';
  public static readonly INIT = 'init';
  public static readonly PHASE = 'phase';
  public static readonly PULL = 'pull';
  public static readonly RUN = 'run';
  public static readonly SECRET = 'secret';
  public static readonly VALIDATE_PHASE = 'validatePhase';
  public static readonly YAML_WRITE = 'yamlWrite';

  public static readonly GET_DOCUMENT_BY_SELECTOR = 'getDocumentsBySelector';
  public static readonly GET_EXECUTOR_DOC = 'getExecutorDoc';
  public static readonly GET_PHASE = 'getPhase';
  public static readonly GET_PHASE_TREE = 'getPhaseTree';
  public static readonly GET_TARGET = 'getTarget';
  public static readonly GET_YAML = 'getYaml';

  public static readonly GET_AIRSHIP_CONFIG_PATH = 'getAirshipConfigPath';
  public static readonly GET_CURRENT_CONTEXT = 'getCurrentContext';
  public static readonly GET_CONTEXTS = 'getContexts';
  public static readonly GET_ENCRYPTION_CONFIGS = 'getEncryptionConfigs';
  public static readonly GET_MANIFESTS = 'getManifests';
  public static readonly GET_MANAGEMENT_CONFIGS = 'getManagementConfigs';

  public static readonly SET_AIRSHIP_CONFIG = 'setAirshipConfig';
  public static readonly SET_CONTEXT = 'setContext';
  public static readonly SET_MANIFEST = 'setManifest';
  public static readonly SET_MANAGEMENT_CONFIG = 'setManagementConfig';
  public static readonly SET_ENCRYPTION_CONFIG = 'setEncryptionConfig';
  public static readonly USE_CONTEXT = 'useContext';

  // UI constants
  public static readonly ANY = 'any';
  public static readonly APPROVED = 'approved';
  public static readonly AUTH = 'auth';
  public static readonly AUTHENTICATE = 'authenticate';
  public static readonly DENIED = 'denied';
  public static readonly HISTORY = 'history';
  public static readonly INITIALIZE = 'initialize';
  public static readonly KEEPALIVE = 'keepalive';
  public static readonly LOG = 'log';
  public static readonly LOGIN = 'login';
  public static readonly REFRESH = 'refresh';
  public static readonly UI = 'ui';
  public static readonly VALIDATE = 'validate';

  public static readonly TASK = 'task';
  public static readonly TASK_END = 'taskEnd';
  public static readonly TASK_START = 'taskStart';
  public static readonly TASK_UPDATE = 'taskUpdate';
}

// Dashboard has the urls of the links that will pop out new dashboard tabs on the left hand side
export class Dashboard {
  name: string;
  baseURL: string;
  path: string;
  isProxied: boolean;
}

// AuthMessage is used to send and auth request and hold the token if it's authenticated
export class Authentication {
  id: string;
  password: string;

  constructor(id?: string | null, password?: string | null) {
    this.id = id;
    this.password = password;
  }
}

