<script>
  import Buttons from '$lib/Buttons.svelte';

  let selectedFiles = [];
  
  async function uploadFile() {
    if (!selectedFiles) {
      return;
    }

    const formData = new FormData();

	for (let i = 0; i < selectedFiles.length; i++) {
    	formData.append('image', selectedFiles[i]);
	}

    try {
      	await fetch('http://localhost:9090/api/upload', {
        method: 'POST',
        body: formData
      });
    } catch (error) {
      console.error(error);
    }
  }

  function handleFileSelect(event) {
    selectedFiles = event.target.files;
  }
</script>

<div style="display: flex;">
<Buttons style="flex: 1;"/>

<div style="flex: 1;" class="container">
  <h2>Upload Folder</h2>
  <input type="file" id="fileInputControl" webkitdirectory="true" multiple accept="image/*" on:change={handleFileSelect}>
  <button on:click={uploadFile}>Upload</button>
</div>

</div>




