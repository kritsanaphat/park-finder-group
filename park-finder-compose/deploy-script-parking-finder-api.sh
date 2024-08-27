export WORKDIR=$(pwd)

function run_staging {

    echo "\n############### DEPLOYING STAGING ###############"
    echo "*** REMOVING OLD CONTAINER ***"
    docker rm -f  parking-finder-api 
    docker rm -f  parking-finder-cronjob
    docker rm -f  parking-finder-demo
    docker rm -f  dozzle




    echo "\n############### DEPLOYING API & CRONJOB ###############"
    docker-compose -f ${WORKDIR}/parking-finder-api/docker-compose.yml pull && \
    docker-compose -f ${WORKDIR}/parking-finder-api/docker-compose.yml up -d
    

} 


function run_prod {

    echo "\n############### DEPLOYING PRODUCTION ###############"
    

}


echo "Deploy To Prod (P,p) or Staging (S,s)? [P,S]"
# read input
if [[ $1 == "P" || $1 == "p" ]]
then
        run_prod        
elif [[ $1 == "S" || $1 == "s" ]]
then
        run_staging
fi