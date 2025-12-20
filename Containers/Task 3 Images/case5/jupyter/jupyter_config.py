# Jupyter Notebook Configuration

c = get_config()

# Allow access from any IP
c.NotebookApp.ip = '0.0.0.0'

# Port
c.NotebookApp.port = 8888

# No browser auto-open
c.NotebookApp.open_browser = False

# Allow root (for container)
c.NotebookApp.allow_root = True

# Disable token for easy access (development only!)
c.NotebookApp.token = ''
c.NotebookApp.password = ''

# Notebook directory
c.NotebookApp.notebook_dir = '/home/jupyter/notebooks'

# Allow remote access
c.NotebookApp.allow_remote_access = True
