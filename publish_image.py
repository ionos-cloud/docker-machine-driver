#!/usr/bin/env python3
"""This is a tool through which the docker-machine-driver image is uploaded in the cloud (Dropbox)."""

import os
import pathlib
import sys

import argparse
import dropbox
import git


def dropbox_connect():
    """Create a connection to Dropbox.

    Returns:
        dbx: New instance of Dropbox client.

    """
    dbx = dropbox.Dropbox(ACCESS_TOKEN)
    return dbx


def dropbox_upload_file(local_path, local_file, dropbox_file_path):
    """Upload a file from the local machine to a path in the Dropbox app directory.

    Args:
        local_path (str): The path to the local file.
        local_file (str): The name of the local file.
        dropbox_file_path (str): The path to the file in the Dropbox app directory.

    Example:
        dropbox_upload_file('.', 'test.csv', '/stuff/test.csv')

    Returns:
        meta: The Dropbox file metadata.

    """
    dbx = dropbox_connect()
    local_file_path = pathlib.Path(local_path) / local_file

    with local_file_path.open("rb") as file:
        meta = dbx.files_upload(
            file.read(), dropbox_file_path, mode=dropbox.files.WriteMode("overwrite")
        )

        return meta


def dropbox_get_link(dropbox_file_path):
    """Get a shared link for a Dropbox file path.

    Args:
        dropbox_file_path (str): The path to the file in the Dropbox app directory.

    Returns:
        link: The shared link.

    """
    try:
        dbx = dropbox_connect()
        shared_link_metadata = dbx.sharing_create_shared_link_with_settings(
            dropbox_file_path
        )
        shared_link = shared_link_metadata.url
        # Replace this in order to download the file directly.
        return shared_link.replace("?dl=0", "?dl=1")
    except dropbox.exceptions.ApiError as exception:
        # IMPORTANT: When using the binary URL in Rancher, a cache is made for that specific URL.
        # If the URL does not differ but the binary does, NO CHANGES will occur!
        # When there is already an image with the same name in the cloud, Drobpox doesn't change the
        # shared link and this is a problem, as you can see from the above important note. We want to
        # avoid that, so we need to force the creation of a new shared link.
        if exception.error.is_shared_link_already_exists():
            # Get the current link.
            shared_link_metadata = dbx.sharing_get_shared_links(dropbox_file_path)
            shared_link = shared_link_metadata.links[0].url
            # Remove the current link.
            dbx.sharing_revoke_shared_link(shared_link)
            # Generate a new shared link.
            shared_link_metadata = dbx.sharing_create_shared_link_with_settings(
                dropbox_file_path
            )
            shared_link = shared_link_metadata.url
            # Replace this in order to download the file directly.
            return shared_link.replace("?dl=0", "?dl=1")


# Constants
LOCAL_PATH = "./bin"
LOCAL_FILE = "docker-machine-driver-ionoscloud"
DROPBOX_FILE_PATH = "/docker-machine-driver-ionoscloud"

# Arguments parser logic and descriptions
DESCRIPTION = """This is a tool through which the docker-machine-driver image is uploaded in the cloud (Dropbox).
                 It prints out an URL that can be used to download the image.
                 The binary must be stored in ./bin and must be named 'docker-machine-driver-ionoscloud'.
                 The authentication is done using a token that can be set using an environment variable named ACCESS_TOKEN."""
parser = argparse.ArgumentParser(description=DESCRIPTION)
args = parser.parse_args()

# Access token validation
ACCESS_TOKEN = os.getenv("ACCESS_TOKEN", None)
if not ACCESS_TOKEN:
    print("Please provide a value for the ACCESS_TOKEN env variable.")
    sys.exit(0)

# Add the commit hash to the dropbox path.
repository = git.Repo(search_parent_directories=True)
commit_hash = repository.head.object.hexsha
dropbox_final_file_path = DROPBOX_FILE_PATH + "-" + commit_hash

dropbox_upload_file(LOCAL_PATH, LOCAL_FILE, dropbox_final_file_path)
print("=======================================")
print("DOWNLOAD THE IMAGE USING THE URL BELOW:")
print(dropbox_get_link(dropbox_final_file_path))
print("=======================================")
